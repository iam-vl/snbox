package main

import (
	"net/http"

	"github.com/iam-vl/snbox/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// func (app *application) routes() *http.ServeMux {
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Create a handler function which wraps our notFound() helper, and then
	// assign it as the custom handler for 404 Not Found responses. You can also
	// set a custom handler for 405 Method Not Allowed responses by setting
	// router.MethodNotAllowed in the same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})
	// Convert ui.Files embedded fs to a http.FS type so it works as a http.FileSystem interface
	// Then pass it to http.FileServer to create a file (server) handler
	fileServer := http.FileServer(http.FS(ui.Files))
	// Our statics are now in the static of folder of the embedded fs. We no longer need to strip the prefix.
	// Any requests with `/static/` will now be passed directly to file server.
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// static file server
	// fileserver := http.FileServer(http.Dir("./ui/static/"))
	// router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))

	// Middleware chain that will contain.
	// Must include nosurf (also inherited by protected)
	dynamic := alice.New(app.sessionManager.LoadAndSave, NoSurf, app.Authenticate)
	// Protected middleware chain:
	protectedChain := dynamic.Append(app.RequireAuth)

	// Unprotected routes
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.HandleHome)) // catch-all
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.HandleViewSnippet))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.HandleSignupForm))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.HandleSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.HandleLoginForm))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.HandleLoginPost))
	// Protected routes
	router.Handler(http.MethodGet, "/snippet/create", protectedChain.ThenFunc(app.HandleSnippetForm))
	router.Handler(http.MethodPost, "/snippet/create", protectedChain.ThenFunc(app.HandleCreateSnippet))
	router.Handler(http.MethodPost, "/user/logout", protectedChain.ThenFunc(app.HandleLogoutUser))

	// router.HandlerFunc(http.MethodGet, "/", app.HandleHome) // catch-all
	// router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.HandleViewSnippet)
	// router.HandlerFunc(http.MethodGet, "/snippet/create", app.HandleSnippetForm)
	// router.HandlerFunc(http.MethodPost, "/snippet/create", app.HandleCreateSnippet)
	// router.HandlerFunc(http.MethodPost, "/head", HandleCustomizeHeaders)

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
