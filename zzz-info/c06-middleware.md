# Middleware 


• Sets useful security headers on every HTTP response.
• Logs the requests received by your application.
• Recovers panics so that they are gracefully handled by your application.
• Create composable middleware chains to help manage and organize your middleware.

## How it works 

Chain of ServeHTTP methods being called after one another. Standard pattern: 

```go
func myMiddleware(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        // Todo: execute middleware logic here...
        next.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
}
```

## Patterns  

Middleware -> Servemux -> app handler // All requests (fe logger) - wrap servemux
Servemux -> middleware > app handler // Requests for specific routes (fe auth)

## Example headers 

```
Content-Security-Policy: default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com
Referrer-Policy: origin-when-cross-origin
X-Content-Type: nosniff
X-Frame-Options: deny
X-XSS-Protection: 0
```

## Misc CH 6.2 

### Flow of control 

```sh
# Original flow:
SecureHeaders -> servemux -> app handlers
# After flow: 
SecureHeaders -> servemux -> app handlers -> servemux -> secureheaders
```

### Early returns 
Example:  
```go
func myMiddleware(next http.Handler) http.Handler {
    return http.Handler(func(w http.ResponseWriter, r *http.Request) {
        // If user isn't authd, send 403 and return to stop executing chain
        if !isAuthorized(r) {
            w.WriteHeader(http.StatusForbidden)
            return
        }
        // otherwise call next handler
        next.ServeHTTP(w, r)
    })
}
```

### Debug CSP issues

CSP headers - blocked by resources - use browser console. Example: 
```
Content Security Policy: the page's settings blocked the loading of a resource at https://... (ex google fonts)
```

## Logging HTTP requests (6.3)

One more mware:  
```go
func (app *application) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}
// Routes 
// LogRequest <-> SecureHeaders <-> servemux <-> handlers
return app.LogRequest(SecureHeaders(mux))
```

## Panic recovery (6.4)

Panic -> application terminated straight away. 
HTTP server assumes: effect of any panic isolated to the goroutine serving the HTTP request. 
Will log a stack trace unwind  the stag to the affected goroutine and close the http conn. 
what if pany in a handler?

```go
func (app *application) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}
	panic("oops! something went wrong") // deliverate panic
	snippets, err := app.snippets.Latest10()
	if err != nil {
		app.ServerError(w, err)
		return
	}
	data := app.NewTemplateData(r)
	data.Snippets = snippets
	fmt.Printf("Year: %+v\n", data.CurrentYear)
	app.Render(w, http.StatusOK, "home.tmpl", data)

}
```
```sh
$ curl -i http://localhost:1111
curl: (52) Empty reply from server # poor - let's dp a 500
```
```
Ping successful
INFO    2024/04/24 01:36:24 Starting server on port: :1111
INFO    2024/04/24 01:36:52 127.0.0.1:41798 - HTTP/1.1 GET /
ERROR   2024/04/24 01:36:52 server.go:3411: http: panic serving 127.0.0.1:41798: oops! something went wrong
goroutine 8 [running]:
net/http.(*conn).serve.func1()
        /snap/go/10585/src/net/http/server.go:1898 +0xbe
panic({0x781d60?, 0x89cbb0?})
        /snap/go/10585/src/runtime/panic.go:770 +0x132
...
```