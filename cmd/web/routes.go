package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	// os_port := os.Getenv("SBOX_PORT")

	// static file server
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("/book", HandleDownloader)

	mux.HandleFunc("/", app.HandleHome) // catch-all
	mux.HandleFunc("/snippet/view", app.HandleViewSnippet)
	mux.HandleFunc("/snippet/create", app.HandleCreateSnippet)
	mux.HandleFunc("/head", HandleCustomizeHeaders)

	return mux

}