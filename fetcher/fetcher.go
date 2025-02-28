package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type FetchConfig struct {
	Headers map[string]string
	Cookies map[string]string
}

func FetchHTML(ctx context.Context, client HttpClient, url string, config ...FetchConfig) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("request creation err: %w", err)
	}

	if len(config) > 0 {
		for key, value := range config[0].Headers {
			req.Header.Add(key, value)
		}

		for key, value := range config[0].Cookies {
			req.AddCookie(&http.Cookie{Name: key, Value: value})
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request execution err: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("response body read err: %w", err)
	}

	return string(body), nil
}
