package fetcher

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type MockHttpClient struct {
	RespBody io.Reader
	Err      error
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if m.RespBody == nil {
		m.RespBody = strings.NewReader("")
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(m.RespBody),
	}, nil
}

type BadReader struct{}

func (b *BadReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Read err")
}

func TestFetchHTML_Success(t *testing.T) {
	ctx := context.Background()
	expectedHTML := "<html><body>Hello</body></html>"
	mockClient := &MockHttpClient{
		RespBody: strings.NewReader(expectedHTML),
	}

	html, err := FetchHTML(ctx, mockClient, "https://www.youtube.com")
	if err != nil {
		t.Fatalf("Err while executed FetchHTML: %v", err)
	}

	if html != expectedHTML {
		t.Errorf("Expected %q, received %q", expectedHTML, html)
	}
}

func TestFetchHTML_RequestError(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockHttpClient{
		Err: errors.New("Request err"),
	}

	_, err := FetchHTML(ctx, mockClient, "https://www.youtube.com")

	if err == nil {
		t.Fatalf("Request err expected but not happened")
	}
}

func TestFetchHTML_RequestCreationError(t *testing.T) {
	mockClient := &MockHttpClient{
		RespBody: strings.NewReader("<html></html>"),
	}

	_, err := FetchHTML(context.Background(), mockClient, "%invalid-url")

	if err == nil {
		t.Fatalf("Request creation err expected but not happened")
	}
}

func TestFetchHTML_BodyReadError(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockHttpClient{
		RespBody: &BadReader{},
	}

	_, err := FetchHTML(ctx, mockClient, "https://www.youtube.com")

	if err == nil {
		t.Fatalf("Response body read err expected but not happened")
	}
}
