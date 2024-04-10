package main

import (
	"log"
	"net/http"
)

const (
	PORT = ":1111"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleHome) // catch-all
	mux.HandleFunc("/snippet/view", HandleViewSnippet)
	mux.HandleFunc("/snippet/create", HandleCreateSnippet)
	mux.HandleFunc("/head", HandleCustomizeHeaders)

	log.Println("Starting server on port:", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatal(err)
	}
}
