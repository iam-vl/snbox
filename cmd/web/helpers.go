package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

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
