package url

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/google/uuid"
)

const (
	idLength = 8
	domain   = "https://api-72ey6bex.nw.gateway.dev"
)

type InvalidURLError struct {
	url string
	err error
}

func (e InvalidURLError) Error() string {
	return fmt.Sprintf("url '%s' is invalid: %v", e.url, e.err)
}

type Relation struct {
	Id         string `json:"id"`
	ShortURL   string `json:"shortUrl"`
	LongURL    string `json:"longUrl"`
	MainDomain string `json:"mainDomain"`
}

// ShortenURL generates a shorter URL for a given long URL and embeds
// everything inside a Relation. The short URL is created by generating
// an ID value. Eventually, this ID is used to look up the longer URL in
// the database.
//
// There will be no check whether some ID is already in use.
func ShortenURL(longURL string) (*Relation, error) {
	var (
		u   *url.URL
		err error
	)

	isProtocolMissing := dotSeparated.Match([]byte(longURL))

	// if only the protocol is missing from the address, we append
	// the HTTP protocol string as a fallback
	if isProtocolMissing {
		longURL = "http://" + longURL
	}

	if u, err = isValidURL(longURL); u == nil {
		return nil, InvalidURLError{url: longURL, err: err}
	}

	rel := &Relation{Id: GenerateID(idLength), LongURL: longURL, MainDomain: u.Hostname()}
	rel.ShortURL = domain + "/" + rel.Id

	return rel, nil
}

// dotSeparated is a regular expression that matches anything that has the
// same dot separation as in an IP address (e.g. 127.0.0.1). Any substring between
// two dots has to be at least one character long.
var dotSeparated = regexp.MustCompile("^(?:\\w+\\.)+\\w+(/\\w+)*$")

// isValidURL tries to validate, if a given address string is a valid URL and return
// the instance of url.URL that was created from parsing.
//
// 'valid' in this case, means:
//
// - there has to be a protocol defined at the beginning of the string ('https://' or 'http://')
//
// - the only allowed non-letter characters in the address are dots (.) and slashes (/)
//
// - sub-routes are allowed (e.g. ...something.com/sub)
func isValidURL(address string) (u *url.URL, err error) {
	if u, err = url.ParseRequestURI(address); err != nil {
		return nil, err
	}

	return
}

// GenerateID generates a random string that contains letters and digits.
func GenerateID(length int) string {
	return uuid.New().String()[0:length]
}
