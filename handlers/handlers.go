package handlers

import "net/http"

func HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from snippetbox"))
}

// /snippet/view
func HandleViewSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Specific snippet"))
}

// snippet/create
func HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet"))
}
