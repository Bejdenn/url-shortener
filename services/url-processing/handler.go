package urlprocessing

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Handle(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var longURL string
		if err := json.NewDecoder(r.Body).Decode(&longURL); err != nil {
			http.Error(rw, "short URL is missing", http.StatusBadRequest)
			return
		}

		fmt.Fprint(rw, longURL)

	default:
		http.Error(rw, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}
