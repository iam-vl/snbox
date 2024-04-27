# CH07 Advanced routing 

julienschmidt/httprouter, go-chi/chi, gorilla/mux
chi / 

```
go get github.com/julienschmidt/httprouter@v1
```

```go
router := httprouter.New()
router.HanderFunc(http.MethodGet, '/snippet/view/:id', app.HandleViewSnippet)
```

## Target endpoints

Before | After  | Handler | Info
---|---|---|---
`ANY`  `/` | `GET`  `/` | `HandleHome` | No catch-all anymore
`ANY`  `/snippet/view?id=123` | `GET`  `/snippet/view/:id`  | `HandleViewSnippet` | 
=none= | `GET` `/snippet/create` | `HandleCreateSnippetForm` | Display HTML form. 
`POST` | `/snippet/create` | `HandleCreateSnippet`  | 
`ANY`  `/static/*filepath` | `GET` `/static/` | `http.Fileserver(...)` | Using the `http.Fileserver()` handler + `http.StripPrefix()`. 

## Routes 

```go
router := httprouter.New()

// static file server
fileserver := http.FileServer(http.Dir("./ui/static/"))
router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))
router.HandlerFunc(http.MethodGet, "/", app.HandleHome) // catch-all

router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.HandleViewSnippet)
// router.HandlerFunc(http.MethodGet, "/snippet/create", app.HandleSnippetForm)
router.HandlerFunc(http.MethodPost, "/snippet/create", app.HandleCreateSnippet)
router.HandlerFunc(http.MethodPost, "/head", HandleCustomizeHeaders)
```

Handlers:
```go
func (app *application) HandleHome(w http.ResponseWriter, r *http.Request) {
	// No need to check now => if r.URL.Path != "/" {
	snippets, err := app.snippets.Latest10()
	if err != nil {
		app.ServerError(w, err)
		return
	}
	data := app.NewTemplateData(r)
	data.Snippets = snippets
	app.Render(w, http.StatusOK, "home.tmpl", data)
}
func (app *application) HandleViewSnippet(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context()) // get params from context
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

func (app *application) HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	// No need to check for POST anymore
	// if r.Method != "POST" {
	// 	app.ClientError(w, http.StatusMethodNotAllowed) 
	// 	return
	// }
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 2
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	http.Re
	
}
```

## Method processing 

```sh
$ curl -i -X POST http://localhost:1111/snippet/view/1
HTTP/1.1 405 Method Not Allowed
Allow: GET, OPTIONS
Content-Security-Policy: default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com
Content-Type: text/plain; charset=utf-8
Referrer-Policy: origin-when-cross-origin
X-Content-Type: nosniff
X-Content-Type-Options: nosniff
X-Frame-Options: deny
X-Xss-Protection: 0
Date: Sat, 27 Apr 2024 13:32:30 GMT
Content-Length: 19

Method Not Allowed
```

## 404 Oddity 

