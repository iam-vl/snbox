# snbox

Learning golang based on Alex Edwards' book.

```sh
curl -i -X POST http://localhost:1111/header
curl https://www.alexedwards.net/static/sb-v2.tar.gz | tar -xvz -C ./ui/static
# 
curl -i -H "Range: bytes=100-199" --output - http:/localhost:1111/static/img/logo.png
go run ./cmd/web
go run ./cmd/web -port=":1234" # ports 0...1023 bound
go run ./cmd/web -help
```

## Misc 

```go
// string to int and int 2 string
i, err := strconv.Atoi("-42")
s := strconv.Itoa(-42)
port := flag.String("port", ":1111", "Server port") 
// also flag.Int, flag.Bool, flag.Float64
// flag.Bool() ~ def true
-> go run ./example -flag
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
Partials:
```
{{define "nav"}}
<nav><a href='/'>Home</a></nav>
{{ end }}
```
in base: `{{ template "nav" . }}`


## Static file server

```go
// 2.8 Add static file server to main.go. Path rel to dir root
fileserver := http.FileServer(http.Dir("./ui/static/"))
// register file server as handler and make sure it matches paths
mux.Handle("/static/", http.StripPrefix("/static", fileserver))
```

## Endpoints

Method | URL | Action
---|---|---
`ANY`  | `/` | `Hello`
`ANY`  | `/snippet/view?id=123` | `Displaying...`
`GET`  | `/snippet/create` | Form 
`POST` | `/snippet/create` | Submit form 
`GET`  | `/user/signup` | Signup form (C 11)
`POST` | `/user/signup` | Create account (C 11)
`GET`  | `/user/login` | Login form (C 11)
`POST` | `/user/login` | Log in user (C 11)
`POST` | `/user/logout` | Log out user (C 11)
`ANY`  | `/static/` | Using the `http.Fileserver()` handler + `http.StripPrefix()`. 

## Installing packages (to be updated)

```sh
go get github.com/go-sql-driver/mysql@v1
go get github.com/justinas/alice@v1
go get github.com/julienschmidt/httprouter@v1
go get github.com/go-playground/form/v4@v4
go get github.com/alexedwards/scs/v2@v2
go get github.com/alexedwards/scs/mysqlstore@latest
```