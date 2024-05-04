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
