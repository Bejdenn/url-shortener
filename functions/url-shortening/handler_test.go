package shorturl

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Bejdenn/url-shortener/functions/url-shortening/url"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const testCollection = "test-url-relations"

var (
	handler = URLHandler{}
	db      *Database
)

func TestMain(m *testing.M) {
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

	db = &Database{Instance: client, TargetCollection: testCollection}

	code := m.Run()
	err = deleteCollection(context.Background(), db, 20)
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestHandler(t *testing.T) {
	tests := []string{
		"http://example.com", "https://www.google.de", "https://elearning.uni-regensburg.de/my/",
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		err := req.ParseForm()
		if err != nil {
			t.Fatalf("could not parse request form: %v", err)
		}
		req.PostForm.Add("longUrl", test)

		rec := httptest.NewRecorder()

		handler.Handle(db, rec, req)

		res := rec.Result()
		body, _ := io.ReadAll(res.Body)

		if res.StatusCode != http.StatusOK {
			t.Fatalf("response contained error: %v", errors.New(string(body)))
		}

		rel := &url.Relation{}
		err = json.Unmarshal(body, rel)
		if err != nil {
			t.Fatalf("error while unmarshalling response body: %v", err)
		}

		if len(rel.ShortURL) == 0 {
			t.Error("expected shortURL but got empty string")
		}
	}

	size := 0
	iter := db.Instance.Collection(testCollection).Documents(context.Background())
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			t.Fatalf("error while trying to get document: %v", err)
		}

		size++
	}

	if size != len(tests) {
		t.Errorf("got %d URL, expected %d in remote database", size, len(tests))
	}
}

func TestHandlerFalseHTTPMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/?longUrl=http://example.com", http.NoBody)
	rec := httptest.NewRecorder()

	handler.Handle(db, rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("got = %d, want = %d", res.StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestHandlerNoURLToProcess(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	rec := httptest.NewRecorder()

	handler.Handle(db, rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("got = %d, want = %d", res.StatusCode, http.StatusBadRequest)
	}
}

func TestShortURLAlreadyExists(t *testing.T) {
	rel, err := url.ShortenURL("https://www.example.com")
	if err != nil {
		t.Fatalf("error occured while shortening URL: %v", err)
	}

	_, err = db.Instance.Collection(testCollection).Doc(rel.Id).Set(context.Background(), rel)
	if err != nil {
		t.Fatalf("error occured while persisting relation: %v", err)
	}

	if exists := handler.Exists(db, rel); !exists {
		t.Errorf("ShortURLExists(%v) = %v, expected = %v", rel, exists, true)
	}
}

func deleteCollection(ctx context.Context, db *Database, batchSize int) error {
	ref := db.Instance.Collection(db.TargetCollection)

	for {
		// Get a batch of documents
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to a
		// WriteBatch.
		batch := db.Instance.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}
