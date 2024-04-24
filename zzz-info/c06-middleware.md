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

CSP headers - blocked by resources - use console.

