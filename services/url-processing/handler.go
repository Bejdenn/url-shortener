package urlprocessing

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ProcessingRequest struct {
	LongURL string `json:"long_url"`
}

func Handle(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var procRequest ProcessingRequest
		if err := json.NewDecoder(r.Body).Decode(&procRequest); err != nil {
			http.Error(rw, "short URL is missing", http.StatusBadRequest)
			return
		}

		fmt.Fprint(rw, procRequest.LongURL)

	default:
		http.Error(rw, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}
