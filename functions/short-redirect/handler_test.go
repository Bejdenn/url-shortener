package shortredirect

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/short-redirect/cf687d69", http.NoBody)
	rec := httptest.NewRecorder()

	Handle(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusMovedPermanently {
		t.Fatalf("got statuscode = %v, want = %v", res.StatusCode, http.StatusMovedPermanently)
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
