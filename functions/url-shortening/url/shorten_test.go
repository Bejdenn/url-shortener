package url

import (
	"regexp"
	"testing"
)

// urlRegex is a regular expression that matches any valid shortened URL that a user would receive.
var urlRegex = regexp.MustCompile("^" + regexp.QuoteMeta(domain) + "\\/\\w{8}$")

func TestShortenValidURLs(t *testing.T) {
	tests := []struct {
		url    string
		domain string
	}{
		{url: "https://www.example.com", domain: "www.example.com"},
		{url: "www.example.com", domain: "www.example.com"},
		{url: "google.de", domain: "google.de"},
	}

	for _, tc := range tests {
		rel, err := ShortenURL(tc.url)
		if err != nil {
			t.Fatalf("error while shortening URL: %v", err)
		}

		if !urlRegex.Match([]byte(rel.ShortURL)) {
			t.Errorf("expected a short URL with schema %s, but got %s", urlRegex.String(), rel.ShortURL)
		}

		if rel.MainDomain != tc.domain {
			t.Errorf("ShortenURL: got %s, wanted: %s", rel.MainDomain, tc.domain)
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
