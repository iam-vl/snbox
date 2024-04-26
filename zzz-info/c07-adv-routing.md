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