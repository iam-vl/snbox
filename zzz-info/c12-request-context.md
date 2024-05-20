# Using request context 

Our only current logic for auth'ing a user is checking for "authenticatedUserId" value in the session data. Helpers:
```go
func (app *application) IsAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserId")
}
```
we can check if the value is real... but we query it often. 
Solution: query in middleware and pass that info down to all subsequent handlers. **Request context**.

Plan: 
* What request context is
* When it's appropriate to use it
* How to use it

## How RQ works 

Every http.request has a context.Context object. 
In a web application, a common usecase: pass info between pieces of middleware and other handlers. 

Syntax:
```go
//r *http.Request
ctx := r.Context() // Retrieve exiting context
ctx = context.WithValue(ctx, "isAuth", true) // create a copy of the context, including a new value
r = r.WithContext(ctx) // create a copy of the request 
```
Can be shortened:
```go
ctx = context.WithValue(ctx, "isAuth", true)
r = r.WithContext(ctx) 
```
retrieving a val:
```go
isAuth, ok := r.Context().Value("isAuth").(bool)
if !ok {
    return errors.New("Coudn't convert value to bool")
}
```
To avoid key collisions: create your own type. Example:
```go
// Declare a custom context key type 
type contextKey string
//Create a constant
const isAuthContextKey = contextKey("isAuth")
// ...

// Set the value in the rq context, using the const
ctx := r.Context()
ctx = context.WithValue(ctx, isAuthContextKey, true)
//...  
// Retrive the value from the ctx using the const as key
isAuth, ok := r.Context().Value(isAuthContextKey).(bool)
if !ok {
	return errors.New("cvould not convert the value to bool")
}
```
