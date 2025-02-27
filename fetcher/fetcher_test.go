package fetcher

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchHTML_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<html><body>Hello</body></html>")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	ctx := context.Background()
	expectedHTML := "<html><body>Hello</body></html>"
	mockClient := &http.Client{}

	html, err := FetchHTML(ctx, mockClient, ts.URL)
	if err != nil {
		t.Fatalf("Error during FetchHTML execution: %v", err)
	}

	if html != expectedHTML {
		t.Errorf("Expected %q, got %q", expectedHTML, html)
	}
}

func TestFetchHTML_RequestExecutionError(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			// Use a transport that simulates a network failure
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return nil, fmt.Errorf("simulated network error")
			},
		},
	}

	_, err := FetchHTML(context.Background(), client, "http://example.com")

	if err == nil {
		t.Fatalf("Expected request execution error, but none occurred")
	}

	if !containsErrorMessage(err, "request execution err") {
		t.Errorf("Expected 'request execution err', but got: %v", err)
	}
}

func TestFetchHTML_RequestCreationError(t *testing.T) {
	ctx := context.Background()

	_, err := FetchHTML(ctx, &http.Client{}, "%invalid-url")

	if err == nil {
		t.Fatalf("Expected request creation error, but none occurred")
	}
}

type MockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.DoFunc(req)
}

type BrokenReader struct{}

func (b *BrokenReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}

func containsErrorMessage(err error, substr string) bool {
	return err != nil && contains(err.Error(), substr)
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr
}

func TestFetchHTML_ResponseBodyReadError(t *testing.T) {
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }
    
    ts := httptest.NewServer(http.HandlerFunc(handler))
    defer ts.Close()
    
    mockClient := &MockHTTPClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            return &http.Response{
                StatusCode: http.StatusOK,
                Body:       io.NopCloser(&BrokenReader{}),
            }, nil
        },
    }
    
    _, err := FetchHTML(context.Background(), mockClient, ts.URL)
    
    if err == nil {
        t.Fatalf("Expected response body read error, but none occurred")
    }
    if !containsErrorMessage(err, "response body read err") {
        t.Errorf("Expected 'response body read err', but got: %v", err)
    }
}
