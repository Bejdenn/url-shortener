package shorturl

import (
	"context"
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

func TestShortURLAlreadyExists(t *testing.T) {
	rel, err := ShortenURL("https://www.example.com")
	if err != nil {
		t.Fatalf("error occured while shortening URL: %v", err)
	}

	Proc.TargetCollection = testCollection

	_, err = Proc.Db.Collection(testCollection).Doc(rel.Id).Set(context.Background(), rel)
	if err != nil {
		t.Fatalf("error occured while persisting relation: %v", err)
	}

	if exists := ShortURLExists(rel); !exists {
		t.Errorf("ShortURLExists(%v) = %v, expected = %v", rel, exists, true)
	}
}
