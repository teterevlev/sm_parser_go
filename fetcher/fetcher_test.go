package fetcher

import "testing"

func TestFetchHTML(t *testing.T) {
	url := "https://www.youtube.com"
	html, err := FetchHTML(url)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(html) == 0 {
		t.Fatal("Empty HTML")
	}
}
