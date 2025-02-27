package main

import (
	"fmt"
	"sm_parser_go/fetcher"
)

func main() {
	url := "https://www.youtube.com"
	html, err := fetcher.FetchHTML(url)
	if err != nil {
		fmt.Println("HTML load err:", err)
		return
	}

	fmt.Println("HTML loaded. Size:", len(html))
}
