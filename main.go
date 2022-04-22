package main

import (
	redirect "github.com/Bejdenn/url-shortener/functions/short-redirect"
	shorturl "github.com/Bejdenn/url-shortener/functions/short-url"
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
