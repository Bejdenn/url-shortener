package main

import (
	"github.com/Bejdenn/url-shortener/functions/redirect"
	shorturl "github.com/Bejdenn/url-shortener/functions/url-shortening"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/url-shortening", shorturl.Handle)
	http.HandleFunc("/", redirect.Handle)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
