package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
    
    

	"sm_parser_go/loadjson"
	"sm_parser_go/fetcher"
	"sm_parser_go/parser"
)

func main() {
	client := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := "https://www.youtube.com"

	// Загружаем конфигурацию
	fetchConfig, err := loadjson.LoadJSON("config.json")

	var htmlContent string // Объявляем переменную заранее

	if err != nil {
		fmt.Println("Config load err:", err)
		// Продолжаем без конфигурации
		htmlContent, err = fetcher.FetchHTML(ctx, client, url)
	} else {
		// Используем конфигурацию при запросе
		htmlContent, err = fetcher.FetchHTML(ctx, client, url, *fetchConfig)
	}

	if err != nil {
		fmt.Println("HTML load err:", err)
		return
	}

	fmt.Println("HTML loaded. Size:", len(htmlContent))

	j, err := parser.GetYTVideoStats(htmlContent)
    fmt.Println(string(*j))
}

