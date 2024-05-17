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
type UserModel struct { DB *sql.DB }
func (m *UserModel) Insert(name, email, pwd string) error { return nil }
func (m *UserModel) Auth(email, pwd string) (int, error) { return 0, nil }
func (m *UserModel) Exists(id int) (bool, error) { return false, nil } 
```
Add new users field to the app: 
```go 
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}
func main() {
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
	// create somewhere to hold custom TLS settings
	tlsConfig := &tls.Config{ 
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
	}

	// Custom http server
	srv := &http.Server{
		// 
	}

	// Need to dereference a pointer
	infoLog.Printf("Starting server on port: %s", *port)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
```

## User signup and pwd encryption 

Create the signup form (ui/html/pages/signup.tmpl)
Include user form struct (handlers):

```go
type UserSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
func (app *application) HandleSignupForm(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = UserSignupForm{}
	app.Render(w, http.StatusOK, "signup.tmpl", data)
}
```


