package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/iam-vl/snbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

// Represent the form data and valid errors.
// Need to be exported  to be read by html/template package
type SnippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

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
// func (app *application) HandleViewSnippet(w http.ResponseWriter, r *http.Request) {
// 	params := httprouter.ParamsFromContext(r.Context())
// 	id, err := strconv.Atoi(params.ByName("id"))
// 	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
// 	if err != nil || id < 1 {
// 		// http.NotFound(w, r)
// 		app.NotFound(w) // Title not blank and < 100 chars long. Add a message if so.(w, err)
// 		return
// 	}
// 	data := app.NewTemplateData(r)
// 	app.Render(w, http.StatusOK, "view.tmpl", data)
// }

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

// snippet/create
func (app *application) HandleSnippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = SnippetCreateForm{
		Expires: 365,
	}
	app.Render(w, http.StatusOK, "create.tmpl", data)
}

// Post to /snippet/create - changed the
// func HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
func (app *application) HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	// Will add post content to r.PostForm
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	// We don't need title and content anymore
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	// Create an instanced of SnippetCreateForm: values + empty map for val errors
	form := SnippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	// this will hold any validation errors
	// form.FieldErrors = make(map[string]string)
	// Title not blank and < 100 chars long. Add a message if so.
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be longer than 100 chars"
	}
	if strings.TrimSpace("content") == "" {
		form.FieldErrors["content"] = "The content cannot be blank"
	}
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "The expires val can only be 1, 7, or 365"
	}
	// Before: If any errors, dump them in plainm HTTP response and return from handler
	// After: If any val errors, rediplay the create template, passing the above struct as dynamic data.
	// Using HTTP 422 Unprocessable Entity in the response to indicate valid. error.
	if len(form.FieldErrors) > 0 {
		data := app.NewTemplateData(r)
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		// fmt.Fprint(w, fieldErrors)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
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
