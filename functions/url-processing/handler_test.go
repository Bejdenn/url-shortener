package urlprocessing

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

var urlRegex = regexp.MustCompile("^https://short-url\\.io\\/\\w{8}$")

const testCollection = "test-urlrelations"

func TestMain(m *testing.M) {
	Proc.TargetCollection = testCollection
	code := m.Run()
	err := deleteCollection(context.Background(), Proc.Db, Proc.Db.Collection(testCollection), 20)
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
		req.ParseForm()
		req.PostForm.Add("longUrl", test)

		rec := httptest.NewRecorder()

		Handle(rec, req)

		res := rec.Result()
		body, _ := io.ReadAll(res.Body)

		if res.StatusCode != http.StatusOK {
			t.Fatalf("response contained error: %v", errors.New(string(body)))
		}

		if shortURL := string(body); len(shortURL) == 0 {
			t.Error("expected shortURL but got empty string")
		}
	}
}

func TestHandlerFalseHTTPMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/?longUrl=http://example.com", http.NoBody)
	rec := httptest.NewRecorder()

	Handle(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("got = %d, want = %d", res.StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestHandlerNoURLToProcess(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	rec := httptest.NewRecorder()

	Handle(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("got = %d, want = %d", res.StatusCode, http.StatusBadRequest)
	}
}

func TestShortenValidURLs(t *testing.T) {
	tests := []string{
		"https://www.example.com", "www.example.com", "google.de",
	}

	for _, v := range tests {
		rel, err := ShortenURL(v)
		if err != nil {
			t.Fatalf("shortening URL was not possible: %v", err)
		}

		if !urlRegex.Match([]byte(rel.ShortURL)) {
			t.Errorf("expected a short URL with schema %s, but got %s", urlRegex.String(), rel.ShortURL)
		}
	}
}

func TestShortenInvalidURLs(t *testing.T) {
	tests := []string{
		"example", "ftp/127.0.0.1", "something?",
	}

	for _, v := range tests {
		if _, err := ShortenURL(v); err == nil {
			t.Errorf("expected error for long URL '%s', but there was none", v)
		}
	}
}

func deleteCollection(ctx context.Context, client *firestore.Client, ref *firestore.CollectionRef, batchSize int) error {

	for {
		// Get a batch of documents
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to a
		// WriteBatch.
		batch := client.Batch()
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
