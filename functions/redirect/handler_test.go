package redirect

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	handler = Handler{}
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

	db = &Database{Instance: client, TargetCollection: "do-not-modify-redirect-test"}

	code := m.Run()
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestHandler(t *testing.T) {
	source := "/cf687d69"
	req := httptest.NewRequest(http.MethodGet, source, http.NoBody)
	rec := httptest.NewRecorder()

	handler.Handle(db, rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusMovedPermanently {
		t.Fatalf("got statuscode = %v, want = %v", res.StatusCode, http.StatusMovedPermanently)
	}

	_, err := res.Location()
	if err != nil {
		t.Fatalf("could not get redirect location: %v", err)
	}
}

func TestExtractPathParam(t *testing.T) {
	tests := []string{
		"abcdef", "12345",
	}

	for _, tc := range tests {
		if param := extractPathParam("/" + tc); param != tc {
			t.Errorf("got = %s, want = %s", param, tc)
		}
	}
}
