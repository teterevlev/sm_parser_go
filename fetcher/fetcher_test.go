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





func TestWithHeaders_EmptyHeaders(t *testing.T) {
	config := FetchConfig{}
	option := WithHeaders(map[string]string{})
	option(&config)

	require.Empty(t, config.Headers, "Headers should be empty")
}

func TestWithHeaders_SingleHeader(t *testing.T) {
	config := FetchConfig{}
	headers := map[string]string{"Content-Type": "application/json"}
	option := WithHeaders(headers)
	option(&config)

	require.Equal(t, headers, config.Headers, "Headers do not match")
}




func TestWithCookies_EmptyCookies(t *testing.T) {
	config := FetchConfig{}
	option := WithCookies(map[string]string{})
	option(&config)

	require.Empty(t, config.Cookies, "Cookies should be empty")
}
func TestWithCookies_SingleCookie(t *testing.T) {
	config := FetchConfig{}
	headers := map[string]string{"session_id": "abc123"}
	option := WithCookies(headers)
	option(&config)

	require.Equal(t, headers, config.Cookies, "Cookies do not match")
}
