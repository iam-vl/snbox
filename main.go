package main

import (
	"log"
	"net/http"
)

const (
	PORT = ":1111"
)

func main() {
	// fmt.Println("hi there")
	// Initialize a new servemux
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleHome)
	log.Println("Starting server on port:", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatal(err)
	}

}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from snippetbox"))
}
