package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/iam-vl/snbox/internal/models"
	"github.com/iam-vl/snbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// Represent the form data and valid errors.
// Need to be exported  to be read by html/template package
// Form parser 8.6: Include struct tags that tell the decoder
// how to map HTML form values to the struct fields
// Ex: we tell the decoder to store the val with name "title" in Title
// "-" - Ignore field during decoding
type SnippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
	// FieldErrors map[string]string
}

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
	//  Retrieve the flash value from the context
	// flash := app.sessionManager.PopString(r.Context(), "flash")
	data := app.NewTemplateData(r)
	data.Snippet = snippet
	// Pass flash to the template
	// data.Flash = flash
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
func (app *application) HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {

	var form SnippetCreateForm
	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	// Create an instanced of SnippetCreateForm: values + empty map for val errors
	// form := SnippetCreateForm{
	// 	Title:   r.PostForm.Get("title"),
	// 	Content: r.PostForm.Get("content"),
	// 	Expires: expires,
	// 	// FieldErrors: map[string]string{},
	// }

	// this will hold any validation errors
	// form.FieldErrors = make(map[string]string)
	// Title not blank and < 100 chars long. Add a message if so.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be longer than 100 chars")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

	if !form.Valid8() {
		// if len(form.FieldErrors) > 0 {
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
	// Add values to the sesh data
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type UserSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) HandleSignupForm(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = UserSignupForm{}
	app.Render(w, http.StatusOK, "signup.tmpl", data)
}
func (app *application) HandleSignupPost(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Create a user")
	var form UserSignupForm
	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 chars long")
	if !form.Valid8() {
		data := app.NewTemplateData(r)
		data.Form = form
		fmt.Printf("Form: %+v\n", form)
		// 422 Unprocessable Content
		app.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}
	// Try creating user rec in db
	fmt.Printf("Creds: %s, %s, %s\n", form.Name, form.Email, form.Password)
	fmt.Println("Insert user model 11")

	err = app.users.Insert(form.Name, form.Email, form.Password)
	fmt.Println("Insert user model 12")
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			fmt.Println("Insert user model 13")

			form.AddFieldError("email", "Email address already in use")
			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			fmt.Println("Insert user model 14")
			app.ServerError(w, err)
		}
		return
	}
	// Otherwise, confirm the operation, and redirect to the login page
	app.sessionManager.Put(r.Context(), "flash", "Your signup has been successul. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther) // HTTP 303
}

type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) HandleLoginForm(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = UserLoginForm{}
	app.Render(w, http.StatusOK, "login.tmpl", data)
}
func (app *application) HandleLoginPost(w http.ResponseWriter, r *http.Request) {
	var form UserLoginForm
	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be a valid address")
	if !form.Valid8() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
	}
	// check if the creds are valid
	id, err := app.users.Auth(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCreds) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.ServerError(w, err)
		}
		return
	}
	// Generate a new session ID when the auth status and priovilege level change
	// For example, if user login / logout.
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// Add the ID of current user to session, so they are now logged in.
	app.sessionManager.Put(r.Context(), "authenticatedUserId", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther) // 403

	// fmt.Fprintln(w, "Auth a user")
}

func (app *application) HandleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// remove the user id from the session data
	app.sessionManager.Remove(r.Context(), "authenticatedUserId")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

// Post to /snippet/create - changed the
// func HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {

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
