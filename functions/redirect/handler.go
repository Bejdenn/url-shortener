package redirect

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var dest404 = `<!DOCTYPE html>
<html lang='en'>
<head>
  <meta charset='UTF-8'>
  <title>Page not found</title>
</head>
<body>
<h1>404 - page not found</h1>
<p>This is a 404 error, which means you've clicked on a bad link or entered an invalid URL.</p>
</body>
</html>`

var template404 = template.Must(template.New("certificate_template").Parse(dest404))

type Handler struct {
}

type Database struct {
	Instance         *firestore.Client
	TargetCollection string
}

func (h *Handler) Handle(db *Database, rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Printf("Registered redirect request for path: %s\n", r.URL.Path)
		id := extractPathParam(r.URL.Path)
		log.Printf("Extracted URL ID from %s: %s\n", r.URL.Path, id)

		iter := db.Instance.Collection(db.TargetCollection).Where("Id", "==", id).Documents(context.Background())
		for {
			doc, err := iter.Next()

			// if there is no URL with the ID, we show a 404 page
			if err == iterator.Done {
				log.Printf("No long URL could be found for ID '%s'\n", id)
				err := template404.Execute(rw, nil)
				if err != nil {
					log.Printf("error: could not execute HTML template: %v", err)
					return
				}

				rw.WriteHeader(http.StatusNotFound)
				return
			}

			if err != nil {
				log.Print(err)
			}

			if longURL, ok := doc.Data()["LongURL"].(string); ok {
				log.Printf("Redirecting successfully to %s\n", longURL)
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
	projectID := "platinum-factor-345219"

	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	db := &Database{Instance: client, TargetCollection: "url-relations"}
	h := Handler{}
	h.Handle(db, rw, r)
}
