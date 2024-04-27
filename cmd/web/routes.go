package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// func (app *application) routes() *http.ServeMux {
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// static file server
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))
	router.HandlerFunc(http.MethodGet, "/", app.HandleHome) // catch-all

	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.HandleViewSnippet)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.HandleSnippetForm)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.HandleCreateSnippet)
	router.HandlerFunc(http.MethodPost, "/head", HandleCustomizeHeaders)

	// mux := http.NewServeMux()
	// // os_port := os.Getenv("SBOX_PORT")
	// mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	// mux.HandleFunc("/book", HandleDownloader)

	// mux.HandleFunc("/", app.HandleHome) // catch-all
	// mux.HandleFunc("/snippet/view", app.HandleViewSnippet)
	// mux.HandleFunc("/snippet/create", app.HandleCreateSnippet)
	// mux.HandleFunc("/head", HandleCustomizeHeaders)

	// LogRequest <-> SecureHeaders <-> servemux <-> handlers
	mwareChain := alice.New(app.RecoverPanic, app.LogRequest, SecureHeaders)
	// return mwareChain.Then(mux)
	return mwareChain.Then(router)
	// return app.RecoverPanic(app.LogRequest(SecureHeaders(mux)))

}
