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

	// static file server
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("/book", HandleDownloader)

	mux.HandleFunc("/", HandleHome) // catch-all
	mux.HandleFunc("/snippet/view", HandleViewSnippet)
	mux.HandleFunc("/snippet/create", HandleCreateSnippet)
	mux.HandleFunc("/head", HandleCustomizeHeaders)

	log.Println("Starting server on port:", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatal(err)
	}
}
