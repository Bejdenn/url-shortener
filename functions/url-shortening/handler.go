package shorturl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Bejdenn/url-shortener/functions/url-shortening/url"
	"google.golang.org/api/iterator"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

const defaultCollection = "url-relations"

type URLHandler struct {
}

type Database struct {
	Instance         *firestore.Client
	TargetCollection string
}

func (h URLHandler) Handle(db *Database, rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Set CORS headers for the main request.
		rw.Header().Set("Access-Control-Allow-Origin", "*")

		err := r.ParseForm()
		if err != nil {
			log.Printf("could not parse request form: %v", err)
			http.Error(rw, "error while processing request", http.StatusInternalServerError)
			return
		}
		longURL := r.PostForm.Get("longUrl")
		if len(longURL) == 0 {
			http.Error(rw, "long URL is missing", http.StatusBadRequest)
			return
		}

		rel, err := url.ShortenURL(longURL)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		for h.Exists(db, rel) {
			rel, _ = url.ShortenURL(longURL)
		}

		_, err = db.Instance.Collection(db.TargetCollection).Doc(rel.Id).Set(context.Background(), rel)
		if err != nil {
			log.Print(err)
			http.Error(rw, "error while processing request", http.StatusInternalServerError)
			return
		}

		payload, err := json.Marshal(rel)
		if err != nil {
			log.Print(err)
			http.Error(rw, "error while processing request", http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")

		if _, err := fmt.Fprint(rw, string(payload)); err != nil {
			log.Printf("error while trying to write response: %v", err)
			http.Error(rw, "error while processing request", http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		rw.Header().Set("Access-Control-Allow-Methods", "POST")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Max-Age", "3600")
		rw.WriteHeader(http.StatusNoContent)
		return

	default:
		http.Error(rw, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func Handle(rw http.ResponseWriter, r *http.Request) {
	// Sets your Google Cloud Platform project ID.
	projectID := "platinum-factor-345219"

	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("could not close Firestore client")
		}
	}(client)

	db := &Database{Instance: client, TargetCollection: defaultCollection}
	h := URLHandler{}
	h.Handle(db, rw, r)
}

// Exists checks whether a Relation is already existing in the database. This is checked by
// comparing the IDs of the relations.
func (h URLHandler) Exists(db *Database, rel *url.Relation) bool {
	tries := 0
	iter := db.Instance.Collection(db.TargetCollection).Where("Id", "==", rel.Id).Documents(context.Background())

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
