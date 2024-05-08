package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

func (app *application) NewTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// dst - tgt destination we wanna decode the form data into
func (app *application) DecodePostForm(r *http.Request, dst any) error {
	// Create Parse form on the request, same way as we did in our createsnippetform handler
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// For all other errors, return them as normal
		return err
	}
	return nil
}

func (app *application) Render(w http.ResponseWriter, status int, page string, tData *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s doesn't exist", page)
		app.ServerError(w, err)
		return
	}
	// Init a new buffer, and write the templ to buffer.
	//  If err, call ServerError()_
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", tData)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// If ok, continue

	w.WriteHeader(status)
	buf.WriteTo(w)
	// err = ts.ExecuteTemplate(w, "base", tData)
	// // err := ts.ExecuteTemplate(w, page, tData)
	// if err != nil {
	// 	app.ServerError(w, err)
	// }
}

// ServerError helper writes and error message and a stack trace to error log
// then sends a generic 500 Internal Server Error response to user
func (app *application) ServerError(w http.ResponseWriter, err error) {
	// Getting a stack trace for the current goroutine and appending it to the message.
	// To see the execuition path of the app via the stack trace is useful when trying to debug errors
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// app.errorLog.Print(trace)
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Same for client errors (like 400 Bad request...)
func (app *application) ClientError(w http.ResponseWriter, status int) {
	// Example: http.StatusText(400) = "Bad Request"
	http.Error(w, http.StatusText(status), status)
}

// Same for 404 not found
func (app *application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}
