package handlers

import "net/http"

func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("hello from snippetbox"))
}

// /snippet/view
func HandleViewSnippet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"name": "Alex"}`))
}

// snippet/create
func HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	// Let's do POST
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
		// VERSION 1
		// w.WriteHeader(405) // you need resp code other than 200 OK
		// w.Write([]byte("Header not allowed"))
		// return
	}
	w.Write([]byte("Create a new snippet"))
}
