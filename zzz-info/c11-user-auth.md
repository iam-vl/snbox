# CH 11 User Authentication

## Plan 

1. Routes setup
2. Creating a user model 
3. User signup and password encryption
4. User login
5. User logout
6. User authorization
7. CSRF protection 

## Route setup

```go
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.HandleSignupForm))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.HandleSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.HandleLoginForm))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.HandleLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.HandleLogoutUser))
```

## User model

```sql
USE snippetbox;
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_pwd CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);
ALTER TABLE users ADD CONSTRAINT users_us_email UNIQUE (email);
```
Add errors (internal/models/errors.go): 
```go
var (
	ErrNoRecord = errors.New("models: no matching records found")
    // New errors
	ErrInvalidCreds = errors.New("models: invalid creds")
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
```
Add model (internal/models/users.go):
```go
type User struct {
	ID         int
	Name       string
	Email      string
	HanshedPwd []byte
	Created    time.Time
}
type UserModel struct {
	DB *sql.DB
}
func (m *UserModel) Insert(name, email, pwd string) error {
	return nil
}
func (m *UserModel) Auth(email, pwd string) (int, error) {
	return 0, nil
}
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
```


