# snbox

Learning golang based on Alex Edwards' book.

```
curl -i -X POST http://localhost:1111/header
curl https://www.alexedwards.net/static/sb-v2.tar.gz | tar -xvz -C ./ui/static
go run ./cmd/web
```

## Template embedding 

```html
{{ define "base" }}
    <title>{{ template "title" . }}</title>
    <main>{{ template "base" . }}</main>
{{ end }}
```

```html
{{ define "title" }}Home{{ end }}
{{ define "main" }}
<h2>Latest snippets </h2>
<p>Nothing to see here yet!</p>
{{ end }}
```



## Static file server

```go
// 2.8 Add static file server to main.go
fileserver := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("/static/", http.StripPrefix("/static", fileserver))
```

## Endpoints

Method | URL | Action
---|---|---
`ANY`  | `/` | `Hello`
`ANY`  | `/snippet/view?id=123` | `Displaying...`
`POST` | `/snippet/create` | `Creating...`
`ANY`  | `/static/` | Using the `http.Fileserver()` handler + `http.StripPrefix()`. 
