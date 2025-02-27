package fetcher

import (
	"fmt"
	"io"
	"net/http"
)

func FetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("request err: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("body read err: %w", err)
	}

	return string(body), nil
}
