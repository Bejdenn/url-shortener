package shortredirect

import (
	"context"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const (
	domain = "https://us-central1-platinum-factor-345219.cloudfunctions.net"
)

var Handler *RedirectHandler

type RedirectHandler struct {
	Db *firestore.Client
}

func init() {
	// Sets your Google Cloud Platform project ID.
	projectID := "platinum-factor-345219"

	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	Handler = &RedirectHandler{Db: client}
}

func (proc *RedirectHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := extractPathParam(r.URL.Path, "/short-redirect")
		iter := Handler.Db.Collection("urlrelations").Where("Id", "==", id).Documents(context.Background())

		for {
			doc, err := iter.Next()

			if err == iterator.Done {
				http.Redirect(rw, r, "https://short-url.io/"+id, http.StatusMovedPermanently)
				break
			}

			if err != nil {
				log.Default().Print(err)
			}

			if longURL, ok := doc.Data()["LongURL"].(string); ok {
				http.Redirect(rw, r, longURL, http.StatusMovedPermanently)
				break

			} else {
				log.Default().Print("error while trying to unmarshall")
			}
		}

	default:
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func extractPathParam(address string, route string) string {
	return strings.TrimPrefix(address, domain+"/"+route+"/")
}

func Handle(rw http.ResponseWriter, r *http.Request) {
	Handler.Handle(rw, r)
}
