package urlprocessing

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Handler(t *testing.T) {
	tests := []struct {
		body string
		want string
	}{
		{body: `{"long_url": "http://example.com"}`, want: "http://example.com"},
	}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "/", strings.NewReader(test.body))
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		Handle(res, req)

		if shortURL := res.Body.String(); shortURL != test.want {
			t.Fatalf("expected shortURL %q but got %q", test.want, shortURL)
		}
	}
}
