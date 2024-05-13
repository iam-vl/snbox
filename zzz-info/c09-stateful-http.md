# CH09 - Stateful HTTP

Plan:
* What session managers are available in Go
* How to use session to safely and securely share data
* How to customize session behavior 

## Choosing session manager

* `gorilla/sessions`
* `alexedwards/scs`

```
go get github.com/alexedwards/scs/v2@v2
go get github.com/alexedwards/scs/mysqlstore@latest
```

## Setting up session manager

```sql
USE snbox;
CREATE TABLE sessions (
    -- Unique random id for each session
    -- actual session data to share between http requests
    -- expiry time for the session (will automatically delete expired ones from the session table)
    token CHAR(43) PRIMARY KEY, 
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
```
Main:
```go
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}
func main() {
    // start
    formDecoder := form.NewDecoder()
	// Configure a sesh manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
    // Rest
}
```
Router:
```go
func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})
	// static file server
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))

    // Changes below
	dynamic := alice.New(app.sessionManager.LoadAndSave)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.HandleHome)) // catch-all
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.HandleViewSnippet))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.HandleSnippetForm))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.HandleCreateSnippet))

	// router.HandlerFunc(http.MethodGet, "/", app.HandleHome) // catch-all
	// router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.HandleViewSnippet)
	// router.HandlerFunc(http.MethodGet, "/snippet/create", app.HandleSnippetForm)
	// router.HandlerFunc(http.MethodPost, "/snippet/create", app.HandleCreateSnippet)
	router.HandlerFunc(http.MethodPost, "/head", HandleCustomizeHeaders)

	mwareChain := alice.New(app.RecoverPanic, app.LogRequest, SecureHeaders)
	return mwareChain.Then(router)
}
```
Version w/out `alice`:
```go
router := httprouter.New()
// ...
router.Handler(http.MethodGet, "/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.HandleHome))) // catch-all
router.Handler(http.MethodGet, "/snippet/view/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.HandleViewSnippet)))
// ...
```

## Flash message 

Handler 1:
```go
func (app *application) HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {

	var form SnippetCreateForm
	// ...
	id, _ := app.snippets.Insert(form.Title, form.Content, form.Expires)
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
```
Template data (templates.go):
```go
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string // Flash message
}
```

Add to the base template:
```html
<main>
	<!-- Display the flash message if any -->
	{{with .Flash}}
		<div class="flash">{{.}}</div>
	{{end}}
	{{ template "main" . }}
</main>
```
Add to handler 2:
```go
func (app *application) HandleViewSnippet(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, _ := strconv.Atoi(params.ByName("id"))
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))

	// Use SnippetModel's Get
	snippet, _ := app.snippets.Get(id)

	//  Retrieve the flash value from the context
	flash := app.sessionManager.PopString(r.Context(), "flash")
	data := app.NewTemplateData(r)
	data.Snippet = snippet
	// Pass flash to the template
	data.Flash = flash
	app.Render(w, http.StatusOK, "view.tmpl", data)
}
```