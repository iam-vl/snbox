package main

import (
	"log"
	"net/http"

	"github.com/iam-vl/snbox/handlers"
)

const (
	PORT = ":1111"
)

func main() {
	// fmt.Println("hi there")
	// Initialize a new servemux
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HandleHome) // catch-all
	mux.HandleFunc("/snippet/view", handlers.HandleViewSnippet)
	mux.HandleFunc("/snippet/create", handlers.HandleCreateSnippet)
	// http.HandleFunc("/", handlers.HandleHome) // doesn't need mux

	log.Println("Starting server on port:", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatal(err)
	}
}
