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