package fetcher

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
    
	"github.com/stretchr/testify/require"
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

    require.Error(t, err)
    require.ErrorContains(t, err, "request execution err")
}

func TestFetchHTML_RequestCreationError(t *testing.T) {
	ctx := context.Background()

	_, err := FetchHTML(ctx, &http.Client{}, "%invalid-url")

	if err == nil {
		t.Fatalf("Expected request creation error, but none occurred")
	}
}
type MockHTTPClient struct{}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(&BrokenReader{}),
	}, nil
}

type BrokenReader struct{}

func (b *BrokenReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}

func TestFetchHTML_ResponseBodyReadError(t *testing.T) {
	mockClient := &MockHTTPClient{}
	
	_, err := FetchHTML(context.Background(), mockClient, "http://example.com")
	
	require.Error(t, err)
	require.ErrorContains(t, err, "response body read err")
}


func TestFetchHTML_WithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get("User-Agent")
		require.Equal(t, "TestAgent", headerValue)
		
		headerValue = r.Header.Get("Accept-Language")
		require.Equal(t, "en-EN", headerValue)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response with headers"))
	}))
	defer server.Close()
	
	client := server.Client()
	config := FetchConfig{
		Headers: map[string]string{
			"User-Agent":      "TestAgent",
			"Accept-Language": "en-EN",
		},
	}
	
	result, err := FetchHTML(context.Background(), client, server.URL, config)
	
	require.NoError(t, err)
	require.Equal(t, "response with headers", result)
}

func TestFetchHTML_WithCookies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		
		foundSessionID := false
		foundUserID := false
		
		for _, cookie := range cookies {
			if cookie.Name == "session_id" && cookie.Value == "abc123" {
				foundSessionID = true
			}
			if cookie.Name == "user_id" && cookie.Value == "12345" {
				foundUserID = true
			}
		}
		
		require.True(t, foundSessionID, "session_id cookie not found or has wrong value")
		require.True(t, foundUserID, "user_id cookie not found or has wrong value")
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response with cookies"))
	}))
	defer server.Close()
	
	client := server.Client()
	config := FetchConfig{
		Cookies: map[string]string{
			"session_id": "abc123",
			"user_id":    "12345",
		},
	}
	
	result, err := FetchHTML(context.Background(), client, server.URL, config)
	
	require.NoError(t, err)
	require.Equal(t, "response with cookies", result)
}

func TestFetchHTML_WithHeadersAndCookies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "TestAgent", r.Header.Get("User-Agent"))
		
		cookies := r.Cookies()
		foundSessionID := false
		
		for _, cookie := range cookies {
			if cookie.Name == "session_id" && cookie.Value == "abc123" {
				foundSessionID = true
				break
			}
		}
		
		require.True(t, foundSessionID, "session_id cookie not found or has wrong value")
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response with headers and cookies"))
	}))
	defer server.Close()
	
	client := server.Client()
	config := FetchConfig{
		Headers: map[string]string{
			"User-Agent": "TestAgent",
		},
		Cookies: map[string]string{
			"session_id": "abc123",
		},
	}
	
	result, err := FetchHTML(context.Background(), client, server.URL, config)
	
	require.NoError(t, err)
	require.Equal(t, "response with headers and cookies", result)
}

func TestFetchHTML_WithEmptyConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response with empty config"))
	}))
	defer server.Close()
	
	client := server.Client()
	config := FetchConfig{
		Headers: map[string]string{},
		Cookies: map[string]string{},
	}
	
	result, err := FetchHTML(context.Background(), client, server.URL, config)
	
	require.NoError(t, err)
	require.Equal(t, "response with empty config", result)
}