package urlprocessing

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

type InvalidURLError struct {
	url string
	err error
}

func (e InvalidURLError) Error() string {
	return fmt.Sprintf("url '%s' is invalid: %v", e.url, e.err)
}

// dotSeparated is a regular expression that matches anything that has the
// same dot separation as in a IP address (e.g. 127.0.0.1). Any part between
// two dots has to be atleast one character long.
var dotSeparated = regexp.MustCompile("^(?:\\w+\\.)+\\w+(\\/\\w+)*$")

type URLRelation struct {
	Id       string `json:"id"`
	ShortURL string `json:"shortUrl"`
	LongURL  string `json:"longUrl"`
}

// ShortenURL generates a shorter URL for a given long URL and embeds
// everything inside a URLRelation. The short URL is created by generating
// a ID value that is used as a path parameter to redirect later.
//
// There will be no check whether some ID is already in use.
func ShortenURL(longURL string) (*URLRelation, error) {
	if valid, err := isValidURL(longURL); !valid {
		return nil, InvalidURLError{url: longURL, err: err}
	}

	rel := &URLRelation{Id: GenerateID(IDLength), LongURL: longURL}
	rel.ShortURL = "https://" + Domain + "/" + rel.Id

	return rel, nil
}

func isValidURL(address string) (bool, error) {
	// try parsing once to filter out the common invalid URLs
	if _, err := url.ParseRequestURI(address); err != nil {
		isProtocolMissing := dotSeparated.Match([]byte(address))

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

func ShortURLExists(rel *URLRelation) bool {
	tries := 0
	iter := Proc.Db.Collection(Proc.TargetCollection).Where("Id", "==", rel.Id).Documents(context.Background())

	for {
		_, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Default().Print(err)
		}

		tries++
	}

	return tries != 0
}
