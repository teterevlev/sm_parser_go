package main


import (
	"context"
	"fmt"
	"net/http"
	"time"
    
	"sm_parser_go/fetcher"
)

func main() {
	client := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := "https://www.youtube.com"
	html, err := fetcher.FetchHTML(ctx, client, url)
	if err != nil {
		fmt.Println("HTML load err:", err)
		return
	}

	fmt.Println("HTML loaded. Size:", len(html))
}
