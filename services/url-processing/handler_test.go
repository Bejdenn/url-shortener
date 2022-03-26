package urlprocessing

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Handler(t *testing.T) {
	longURL := "https://example.com"
	req := httptest.NewRequest("POST", "/url-processing", strings.NewReader(longURL))
	res := httptest.NewRecorder()

	Handle(res, req)

	if shortURL := res.Body.String(); shortURL != longURL {
		t.Fatalf("expected shortURL %q but got %q", longURL, shortURL)
	}
}
