package urlprocessing

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

var (
	urlRegex = regexp.MustCompile("^https://short-url\\.io\\/\\w{8}$")
)

func TestHandler(t *testing.T) {
	tests := []string{
		"http://example.com", "https://www.google.de", "https://elearning.uni-regensburg.de/my/",
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodGet, "/?longUrl="+test, http.NoBody)
		rec := httptest.NewRecorder()

		Handle(rec, req)

		res := rec.Result()
		body, _ := io.ReadAll(res.Body)

		if res.StatusCode != http.StatusOK {
			t.Fatalf("response contained error: %v", errors.New(string(body)))
		}

		if shortURL := string(body); len(shortURL) == 0 {
			t.Error("expected shortURL but got empty string")
		}
	}
}

func TestShortenValidURLs(t *testing.T) {
	tests := []string{
		"https://www.example.com", "www.example.com", "google.de",
	}

	for _, v := range tests {
		shortURL, err := ShortenURL(v)
		if err != nil {
			t.Fatalf("shortening URL was not possible: %v", err)
		}

		if !urlRegex.Match([]byte(shortURL)) {
			t.Errorf("expected a short URL with schema %s, but got %s", urlRegex.String(), shortURL)
		}
	}
}

func TestShortenInvalidURLs(t *testing.T) {
	tests := []string{
		"example", "ftp/127.0.0.1",
	}

	for _, v := range tests {
		if _, err := ShortenURL(v); err == nil {
			t.Errorf("expected error for long URL '%s', but there was none", v)
		}
	}
}
