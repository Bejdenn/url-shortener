package shortredirect

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const (
	domain  = "https://us-central1-platinum-factor-345219.cloudfunctions.net"
	route   = "/short-redirect"
	dest404 = "https://short-url.io/"
)

var handler *RedirectHandler

type RedirectHandler struct {
	database *firestore.Client
}

func init() {
	projectID := "platinum-factor-345219"

	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	handler = &RedirectHandler{database: client}
}

func (h *RedirectHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Printf("Registered redirect request for path: %s\n", r.URL.Path)
		id := extractPathParam(r.URL.Path)
		fmt.Printf("Extracted URL ID from %s: %s\n", r.URL.Path, id)

		iter := handler.database.Collection("urlrelations").Where("Id", "==", id).Documents(context.Background())
		for {
			doc, err := iter.Next()

			// if there is no URL with the ID and we redirect to a 404 page
			if err == iterator.Done {
				log.Printf("No long URL could be found for ID '%s'\n", id)
				http.Redirect(rw, r, dest404+id, http.StatusMovedPermanently)
				return
			}

			if err != nil {
				log.Print(err)
			}

			if longURL, ok := doc.Data()["LongURL"].(string); ok {
				fmt.Printf("Redirecting successfully to %s\n", longURL)
				http.Redirect(rw, r, longURL, http.StatusMovedPermanently)
				break

			} else {
				log.Print("error while trying to unmarshall")
			}
		}

	default:
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func extractPathParam(address string) string {
	return strings.ReplaceAll(address, "/", "")
}

func Handle(rw http.ResponseWriter, r *http.Request) {
	handler.Handle(rw, r)
}
