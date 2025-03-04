package loadjson

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"sm_parser_go/fetcher"
)

// ConfigData структура для хранения конфигурации
type ConfigData struct {
	CookiesYT map[string]string `json:"cookies_yt"`
	HeadersYT map[string]string `json:"headers_yt"`
	// Другие поля, которые мы пока игнорируем
}

// LoadJSON загружает конфигурацию из json файла
func LoadJSON(filePath string) (*fetcher.FetchConfig, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var configData ConfigData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&configData); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	fetchConfig := &fetcher.FetchConfig{
		Cookies: configData.CookiesYT,
		Headers: configData.HeadersYT,
	}

	return fetchConfig, nil
}