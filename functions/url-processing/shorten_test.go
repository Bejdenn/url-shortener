package urlprocessing

import (
	"regexp"
	"testing"
)

var urlRegex = regexp.MustCompile("^https://" + regexp.QuoteMeta(Domain) + "\\/\\w{8}$")

func TestShortenValidURLs(t *testing.T) {
	tests := []string{
		"https://www.example.com", "www.example.com", "google.de",
	}

	for _, tc := range tests {
		rel, err := ShortenURL(tc)
		if err != nil {
			t.Fatalf("error while shortening URL: %v", err)
		}

		if !urlRegex.Match([]byte(rel.ShortURL)) {
			t.Errorf("expected a short URL with schema %s, but got %s", urlRegex.String(), rel.ShortURL)
		}
	}
}

func TestShortenInvalidURLs(t *testing.T) {
	tests := []string{
		"example", "ftp/127.0.0.1", "something?",
	}

	for _, tc := range tests {
		if _, err := ShortenURL(tc); err == nil {
			t.Errorf("expected error for URL '%s', but there was none", tc)
		}
	}
}
