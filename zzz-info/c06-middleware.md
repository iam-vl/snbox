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

Middleware -> Servemux -> app handler // All requests (fe logger)
Servemux -> middleware > app handler // Requests for specific routes (fe auth)