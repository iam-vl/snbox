# Notes 

## Disable directory list

Add a blank index.html to each dir 
```sh
$ find ./ui/static -type d -exec touch {}/index.html \;
```
A more complicated (but arguably better) solution is to create a custom implementation of `http.FileSystem`, and have it return an `os.ErrNotExist` error for any directories. A full explanation and sample code can be found here: [How to Disable FileServer Directory Listings](https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings).

## About http.Handler interface

```go
type Handler interface { 
    ServeHTTP(ResponseWriter, *Request) 
}
```
Simplest handler: 
```go
type HandleHome struct {}
func (h *HandleHome) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]bytes("This is my home page"))
}
```
This can lead to:
```go
mux := http.NewServeMux()
mux.Handle("/", &HandleHome{})
```

```go
mux.HandleFunc("/", home) = mux.Handle("/", http.HandlerFunc(home))
```