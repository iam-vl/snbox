# CH08 Process Forms 

## Create and process basic forms

```go
func (app *application) HandleShowForm(w http.ResponseWriter, r *http.Request) {
    data := app.NewTemplateData(r)
    app.render(w, http.StatusOK, "create.tmpl", data)
}
func (app *application) HandleProcessForm(w http.ResponseWriter, r *http.Request) {

    err := r.ParseForm()
    if err := nil {
        app.ClientError(w, http.StatusBadRequest)
        return 
    }
    // r.PostFort.Get() always returns a string
    title := r.PostForm.Get("title") 
    content := r.PostForm.Get("content")
    expires, err := strconv.Atoi(r.PostForm.Get("expires")) 
    if err = nil {
        app.ClientError(w, http.StatusBadRequest)
        return 
    }
    id, err := app.snippets.Insert(title, content, expires)
    if err = nil {
        app.ServerError(w, err)
        return 
    }
    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
```

## Other forms

r.PostForm | Request body only (POST, PATCH, PUT) 
r.Form() | all methods, all vals. But can result in name conflict
r.ParseForm().Get() -> r.PostFormValue() | Returns first val only
r.Form.Get() -> r.FormValue()

## Multiple value fields

```html
<input type="checkbox" name="items" value="foo">Foo
<input type="checkbox" name="items" value="bar">Bar
<input type="checkbox" name="items" value="baz">Baz
```
```go
for i, item := range r.PostForm["items"] {
    fmt.Fprintf(w, "%d: Item %s\n", i, item)
}
```

## Limiting form size 

By def, max form size 10MB. To change it:
```go
r.Body = http.MaxBytesReader(w, r.Body, 4096) // 4096 Bytes max
r.ParseForem will try to read 4096. If more, error
err := r.ParseForm()
if err != nil {
    http.Error(w, "Bad Request", http.StatusBadRequest)
    return
}
```

## Validate form 

```go
func (app *application) HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	// Parse the form an get values
    // ... 
	// Hold validation errors
	fieldErrors := make(map[string]string)
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 { // "Zoë" Unicode len("Zoë") will be 4
		fieldErrors["title"] = "This field cannot be longer than 100 chars"
	}
	if strings.TrimSpace("content") == "" {
		fieldErrors["content"] = "The content cannot be blank"
	}
	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "The expires val can only be 1, 7, or 365"
	}
	// If any errors, dump them in plainm HTTP response and return from handler
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}
    // Go on as usual
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
```


## Parsing forms automatically

To automatically decode the form data into `CreateSnippetForm`: 
```
go get github.com/go-playground/form/v4@v4
```

```go
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	// Set up the basics 
	formDecoder := form.NewDecoder()
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}
    // Start a server ...
}
```
Handlers:
```go
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
}
```
Handler function:
```go
	// expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	// if err != nil {
	// 	app.ClientError(w, http.StatusBadRequest)
	// 	return
	// }
    var form SnippetCreateForm
	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest) // return a 400
		return
	}
	// Create an instanced of SnippetCreateForm: values + empty map for val errors
	// form := SnippetCreateForm{
	// 	Title:   r.PostForm.Get("title"),
	// 	Content: r.PostForm.Get("content"),
	// 	Expires: expires,
	// 	// FieldErrors: map[string]string{},
	// }
```

>**Issue:**
>The `app.formDecoder.Decode()` requires a non-nil pointer as the target decode destination. If a nil-pointer, it'll return a `form.InvalidDecoderError`. 
>We need to manage this as a special case.  

Solution: create a new `DecodePostForm()` helper:
* Calls `r.ParseForm()` on the current request. 
* Calls `app.formDecoder.Decode()` to unpack the form data to a TGT destination. 
* Checks for a `form.InvalidDecoderError` error and triggers a panic if we ever see it. 

Updating `helpers.go`:
```go
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
```


Updating handlers:
```go
func (app *application) HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	// Will add post content to r.PostForm
	// err := r.ParseForm()
	// if err != nil {
	// 	app.ClientError(w, http.StatusBadRequest)
	// 	return
	// }
	var form SnippetCreateForm
	// err = app.formDecoder.Decode(&form, r.PostForm)\
	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
    // .. 
}
```