# Displaying Errrors 

## Template data 

Templates.go: 
```go
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
}
```

```go
func (app *application) HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	// Will add post content to r.PostForm
	_ := r.ParseForm()
	// if err != nil { app.ClientError(w, http.StatusBadRequest) return }
	// We don't need title and content anymore
	expires, _ := strconv.Atoi(r.PostForm.Get("expires"))
	// if err != nil { app.ClientError(w, http.StatusBadRequest) return }
	form := SnippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}
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
	// HTTP 422 Unprocessable Entity in the response to indicate valid. error.
	if len(form.FieldErrors) > 0 {
		data := app.NewTemplateData(r)
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		// fmt.Fprint(w, fieldErrors)
		return
	}
	id, _ := app.snippets.Insert(form.Title, form.Content, form.Expires)
	// if err != nil { app.ServerError(w, err) return }
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
```

```html
{{ define "main" }}
    <form action="/snippet/create" method="POST">
        <div>
            <label>Title</label>
            <!-- Render val errors -->
            {{ with .Form.FieldErrors.title }}
                <label class="error">{{.}}</label>
            {{ end }}
            <!-- Repopulate the title data by setting 'value' -->
            <input type="text" name="title" value="{{.Form.Title}}">
        </div>
        <div>
            <label>Content</label>
            <!-- Render val errors -->
            {{ with .Form.FieldErrors.content }}
                <label class="error">{{.}}</label>
            {{ end }}
            <input type="text" name="title">
            <!-- Repopulate the content data -->
            <textarea name="content">{{ .Form.Content }}</textarea>
            <!-- <textarea name="content" id="" cols="30" rows="10"></textarea> -->
        </div>
        <div>
            <label>Delete in:</label>
            {{ with .Form.FieldErrors.expires }}
                <label class="error">{{.}}</label>
            {{ end }}
            <!-- Check repopulated expires values -->
            <input type="radio" name="expires" value="365" {{if (eq .Form.Expires 365) }}checked{{end}}>One year
            <input type="radio" name="expires" value="7" {{if (eq .Form.Expires 7) }}checked{{end}}>One week
            <input type="radio" name="expires" value="1" {{if (eq .Form.Expires 1) }}checked{{end}}>One day
        </div>
        <div>
            <input type="submit" value="Publish snippet">
        </div>
    </form>
{{ end }}
```