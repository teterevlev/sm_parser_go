package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
    
	"sm_parser_go/loadjson"
	"sm_parser_go/fetcher"
)

func main() {
	client := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	url := "https://www.youtube.com"
	
	// Загружаем конфигурацию
	fetchConfig, err := loadjson.LoadJSON("config.json")
	if err != nil {
		fmt.Println("Config load err:", err)
		// Продолжаем без конфигурации если не удалось загрузить
		html, err := fetcher.FetchHTML(ctx, client, url)
		if err != nil {
			fmt.Println("HTML load err:", err)
			return
		}
		fmt.Println("HTML loaded without config. Size:", len(html))
		return
	}
	
	// Используем конфигурацию при запросе
	html, err := fetcher.FetchHTML(ctx, client, url, *fetchConfig)
	if err != nil {
		fmt.Println("HTML load err:", err)
		return
	}
	fmt.Println("HTML loaded with config. Size:", len(html))
}