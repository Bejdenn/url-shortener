package urlprocessing

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

var Proc *URLProcessing

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/Users/dennisbejze/Projects/url-shortener/service_account.json")

	// Sets your Google Cloud Platform project ID.
	projectID := "platinum-factor-345219"

	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	Proc = &URLProcessing{Db: client, TargetCollection: DefaultCollection}
}

type InvalidURLError struct {
	url string
	err error
}

func (e InvalidURLError) Error() string {
	return fmt.Sprintf("url '%s' is invalid: %v", e.url, e.err)
}

const (
	Domain            = "https://short-url.io"
	IDLength          = 8
	DefaultCollection = "urlrelations"
)

// dotSeparated is a regular expression that matches anything that has the
// same dot separation as in a IP address (e.g. 127.0.0.1). Any part between
// two dots has to be atleast one character long.
var dotSeparated = regexp.MustCompile("^(?:\\w+\\.)+\\w+(\\/\\w+)*$")

type URLProcessing struct {
	Db               *firestore.Client
	TargetCollection string
}

func (proc *URLProcessing) Handle(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		longURL := r.PostForm.Get("longUrl")
		if len(longURL) == 0 {
			http.Error(rw, "long URL is missing", http.StatusBadRequest)
			return
		}

		rel, err := ShortenURL(longURL)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = proc.Db.Collection(proc.TargetCollection).Doc(rel.Id).Set(context.Background(), rel)
		if err != nil {
			log.Default().Print(err)
		}

	default:
		http.Error(rw, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func Handle(rw http.ResponseWriter, r *http.Request) {
	Proc.Handle(rw, r)
}

type URLRelation struct {
	Id       string
	ShortURL string
	LongURL  string
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
	rel.ShortURL = Domain + "/" + rel.Id

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
