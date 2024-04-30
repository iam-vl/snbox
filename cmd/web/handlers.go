package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/iam-vl/snbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) HandleHome(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	app.NotFound(w)
	// 	return
	// }
	// panic("oops! something went wrong") // deliverate panic
	snippets, err := app.snippets.Latest10()
	if err != nil {
		app.ServerError(w, err)
		return
	}
	data := app.NewTemplateData(r)
	data.Snippets = snippets
	// Use render helper
	fmt.Printf("Year: %+v\n", data.CurrentYear)
	app.Render(w, http.StatusOK, "home.tmpl", data)
}

// /snippet/view?id=123
func (app *application) HandleViewSnippet(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.NotFound(w)
		return
	}
	// Use SnippetModel's Get
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(w)
		} else {
			app.ServerError(w, err)
		}
		return
	}
	data := app.NewTemplateData(r)
	data.Snippet = snippet
	app.Render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) HandleSnippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	app.Render(w, http.StatusOK, "create.tmpl", data)
}

// snippet/create - changed the
// func HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
func (app *application) HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	// No need to check for POst anymore
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 2
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// w.Write([]byte("Creating a new snippet"))
	// http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// snippet/create
func HandleCustomizeHeaders(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=31536000") // overrides
	w.Header().Add("Cache-Control", "public")
	w.Header().Add("Cache-Control", "max-age=31536000")
	// Avoid canonicalization
	// w.Header()["X-XSS-Protection"] = []string("1;mode=block")
	fmt.Printf("Header before deleting / suppressing: %+v\n", w.Header())
	fmt.Printf("Date before suppressing: %+v\n", w.Header().Get("Date"))
	w.Header()["Date"] = nil // suppress a system generated header
	fmt.Printf("Header before deleting: %+v\n", w.Header())
	fmt.Printf("First val: %+v\n", w.Header().Get("Cache-Control")) // first val
	fmt.Printf("Entire header after deleting: %+v\n", w.Header())
	w.Header().Del("Cache-Control")
	fmt.Println("===========")
	fmt.Printf("Header after deleting: %+v\n", w.Header())

	w.Write([]byte(`{"name": "Alex"}`))
}

func HandleDownloader(w http.ResponseWriter, r *http.Request) {
	fmt.Println("downloading")
	http.ServeFile(w, r, "./us/static/lets-go.epub")
}
