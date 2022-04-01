package shortredirect

import "testing"

func TestExtractPathParam(t *testing.T) {
	tests := []string{
		"abcdef", "12345",
	}

	for _, tc := range tests {
		if param := extractPathParam("/short-redirect/"+tc, "/short-redirect"); param != tc {
			t.Errorf("got = %s, want = %s", param, tc)
		}
	}
}
