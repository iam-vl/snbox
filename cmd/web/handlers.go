package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) HandleHome(w http.ResponseWriter, r *http.Request) {
	// func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// http.NotFound(w, r)
		app.NotFound(w) // using NotFound() helper
		return
	}
	fmt.Println("Starting home...")

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	// ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	if err != nil {
		// log.Panicln(err)
		// http.Error(w, "Error parsing template", 500)
		app.ServerError(w, err)
		return
	}
	// if err := ts.Execute(w, nil); err != nil {
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil { // use base
		// log.Println(err.Error())
		app.ServerError(w, err)
		// app.errorLog.Print(err.Error())
		// http.Error(w, "Internal Server Error", 500)
	}
	// w.Write([]byte("hello from snippetbox"))
}

// /snippet/view?id=123
func (app *application) HandleViewSnippet(w http.ResponseWriter, r *http.Request) {
	// func HandleViewSnippet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("starting view snippet")
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.NotFound(w)
		return
	}
	// http.ResponseWriter (w.Write) <- io.Writer interface
	fmt.Fprintf(w, "Displaying a snippet with ID: %v...\n", id)
	// w.Write([]byte(`{"name": "Alex", "id": %d}`))
}

// snippet/create - changed the
// func HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
func (app *application) HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	// Let's do POST
	fmt.Println("posting")
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		app.ClientError(w, http.StatusMethodNotAllowed) // using ClientError()
		return
		// VERSION 1
		// w.WriteHeader(405) // you need resp code other than 200 OK
		// w.Write([]byte("Header not allowed"))
		// return
	}
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// w.Write([]byte("Creating a new snippet"))
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
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
