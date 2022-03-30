package urlprocessing

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/google/uuid"
)

type InvalidURLError struct {
	url string
	err error
}

func (e InvalidURLError) Error() string {
	return fmt.Sprintf("url '%s' is invalid: %v", e.url, e.err)
}

const (
	Domain   = "https://short-url.io"
	IDLength = 8
)

func Handle(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		longURL := r.URL.Query().Get("longUrl")
		if len(longURL) == 0 {
			http.Error(rw, "long URL is missing", http.StatusBadRequest)
			return
		}

		shortURL, err := ShortenURL(longURL)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprint(rw, shortURL)

	default:
		http.Error(rw, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// ShortenURL generates a shorter URL for a given long URL. The short URL is created by
// creating a hash value that is used as a ID for routing.
//
// There will be no check whether some ID is already in use.
func ShortenURL(longURL string) (string, error) {
	if valid, err := isValidURL(longURL); !valid {
		return "", InvalidURLError{url: longURL, err: err}
	}

	return Domain + "/" + GenerateID(IDLength), nil
}

func isValidURL(address string) (bool, error) {
	// try parsing once to filter out the common invalid URLs
	if _, err := url.ParseRequestURI(address); err != nil {
		isProtocolMissing, err := regexp.Match("^(?:\\w+\\.)+\\w+(\\/\\w+)*$", []byte(address))
		if err != nil {
			return false, err
		}

		// if only the protocol is missing from the address, then we append
		// the HTTP protocol string as a fallback and retry
		if isProtocolMissing {
			address = "http://" + address
		}

		if _, err := url.ParseRequestURI(address); err != nil {
			return false, err
		}
	}

	return true, nil
}

func GenerateID(length int) string {
	return uuid.New().String()[0:length]
}
