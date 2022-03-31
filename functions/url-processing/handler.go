package urlprocessing

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
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

const (
	Domain            = "short-url.io"
	IDLength          = 8
	DefaultCollection = "urlrelations"
)

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

		payload, err := json.Marshal(rel)
		if err != nil {
			log.Default().Print(err)
		}

		fmt.Fprint(rw, string(payload))

	default:
		http.Error(rw, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func Handle(rw http.ResponseWriter, r *http.Request) {
	Proc.Handle(rw, r)
}
