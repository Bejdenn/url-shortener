package urlprocessing

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

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

		rel := &URLRelation{}
		err := json.Unmarshal(body, rel)
		if err != nil {
			t.Fatalf("error while unmarshalling response body: %v", err)
		}

		if len(rel.ShortURL) == 0 {
			t.Error("expected shortURL but got empty string")
		}
	}

	size := 0
	iter := Proc.Db.Collection(testCollection).Documents(context.Background())
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
