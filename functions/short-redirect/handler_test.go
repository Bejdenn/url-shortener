package shortredirect

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	source := domain + "/short-redirect/cf687d69"
	req := httptest.NewRequest(http.MethodGet, source, http.NoBody)
	rec := httptest.NewRecorder()

	Handle(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusMovedPermanently {
		t.Fatalf("got statuscode = %v, want = %v", res.StatusCode, http.StatusMovedPermanently)
	}

	url, err := res.Location()
	if err != nil {
		t.Fatalf("could not get redirect location: %v", err)
	}

	if url.RawPath != source {

	}
}

func TestExtractPathParam(t *testing.T) {
	tests := []string{
		"abcdef", "12345",
	}

	for _, tc := range tests {
		if param := extractPathParam(domain + "/short-redirect/" + tc); param != tc {
			t.Errorf("got = %s, want = %s", param, tc)
		}
	}
}
