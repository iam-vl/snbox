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

## Create UserModel.Exists()

```go
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
```
## Create context key 

context.go
```go
package main
type contextKey string
const isAuthContextKey = contextKey("isAuth")
```

## Create authenticate() middleware

Let's create a new middleware that:
1. Retrieves the user id from session data.
2. Checks the db to see if the id correponds tyo a real user. 
3. Update the reqt context to include the context ley with value `true`. 

Middleware: 
```go
func (app *application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the userid val from the session using GetInt()
		// Will return an int (userid) or zero (none)
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}
		// check to see if a user with this Id exists
		exists, err := app.users.Exists(id)
		if err != nil {
			app.ServerError(w, err)
			return
		}
		// if ok, we create a copy of the request and assign it to r 
		if exists {
			ctx := context.WithValue(r.Context(), isAuthContextKey, true)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
```
Include app.Authenticate in the dynamic middleware chain (routes):
```go
dynamic := alice.New(app.sessionManager.LoadAndSave, NoSurf, app.Authenticate)
```

## Update isAuthenticated() helper 

```go
func (app *application) IsAuthenticated(r *http.Request) bool {
	isAuth, ok := r.Context().Value(isAuthContextKey).(bool)
	if !ok {
		return false
	}
	return isAuth
}
```
