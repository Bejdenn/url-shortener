package url

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/google/uuid"
)

const (
	idLength = 8
	domain   = "https://api-url-shortener-72ey6bex.nw.gateway.dev"
)

type InvalidURLError struct {
	url string
	err error
}

func (e InvalidURLError) Error() string {
	return fmt.Sprintf("url '%s' is invalid: %v", e.url, e.err)
}

type Relation struct {
	Id       string `json:"id"`
	ShortURL string `json:"shortUrl"`
	LongURL  string `json:"longUrl"`
}

// ShortenURL generates a shorter URL for a given long URL and embeds
// everything inside a Relation. The short URL is created by generating
// an ID value. Eventually, this ID is used to look up the longer URL in
// the database.
//
// There will be no check whether some ID is already in use.
func ShortenURL(longURL string) (*Relation, error) {
	if valid, err := isValidURL(longURL); !valid {
		return nil, InvalidURLError{url: longURL, err: err}
	}

	rel := &Relation{Id: GenerateID(idLength), LongURL: longURL}
	rel.ShortURL = domain + "/" + rel.Id

	return rel, nil
}

// dotSeparated is a regular expression that matches anything that has the
// same dot separation as in an IP address (e.g. 127.0.0.1). Any substring between
// two dots has to be at least one character long.
var dotSeparated = regexp.MustCompile("^(?:\\w+\\.)+\\w+(/\\w+)*$")

// isValidURL tries to validate, if a given address string is a valid URL.
//
// 'valid' in this case, means:
//
// - there has to be a protocol defined in the beginning of the string ('https://' or 'http://')
//
// - the only allowed non-letter character in the address is a dot
//
// - sub-routes are allowed (e.g. ...something.com/sub)
func isValidURL(address string) (bool, error) {
	// try parsing once to filter out the common invalid URLs
	if _, err := url.ParseRequestURI(address); err != nil {
		isProtocolMissing := dotSeparated.Match([]byte(address))

		// if only the protocol is missing from the address, we append
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

// GenerateID generates a random string that contains letters and digits.
func GenerateID(length int) string {
	return uuid.New().String()[0:length]
}
