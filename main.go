package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	rwFile, err := os.OpenFile("/Users/lizifeng/code/go/src/pull-zhihu-comment/comment.txt", os.O_RDWR|os.O_APPEND, 0766)
	if err != nil {
		log.Fatal("打开文件发生错误:", err)
	}
	defer rwFile.Close()

	url := "https://www.zhihu.com/api/v4/answers/1713375801/comments?limit=20&status=open&order=obverse&offset="
	count := 0
	for i := 10000; i >= 0; i = i - 20 {
		fmt.Println("requesting ", url+strconv.Itoa(i))
		body := requestUrl(url + strconv.Itoa(i))
		saveResult(body, &count, rwFile)
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
}

func saveResult(body []byte, count *int, file *os.File) {
	result := Result{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}

	for i := len(result.Data) - 1; i >= 0; i-- {
		Comment := result.Data[i]
		if "author" == Comment.Author.Role {
			*count++
			date := time.Unix(Comment.CreatedTime, 0).Format("2006-01-02 15:04:05")
			line := fmt.Sprintf("[%s] %s. %s\n", date, strconv.Itoa(*count), Comment.Content)

			if _, err := file.WriteString(line); err != nil {
				log.Fatalf("wrtie '%s' error: %s", line, err)
			}
		}
	}
}

func requestUrl(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}
	return body
}

type Result struct {
	Data []Comment `json:"data"`
}

type Comment struct {
	Type        string `json:"type"`
	Content     string `json:"content"`
	CreatedTime int64  `json:"created_time"`
	Author      struct {
		Role string `json:"role"`
	} `json:"author"`
}
