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
## Validate input

```go
var EmailRegex = regexp.MustCompile("`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`")

func MinChars(val string, n int) bool {
	// True if a val contains at least n chars
	return utf8.RuneCountInString(val) >= n
}
func Matches(val string, rx *regexp.Regexp) bool {
	// True if the val matches the regex
	return rx.MatchString(val)
}
```
Process the form and run the valid8n tests (handlers):
```go
```

Using bcrypt:
```
go get golang.org/x/crypto/bcrypt@latest
```
```go
func GenerateHash(password string) {
	// Include password and cost (4-31). 12 is reasonable minimum.
	// Returns a 60-char string
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	fmt.Printf("Hash value: %+v\t\nHash type: %T\n", hash, hash)
	// Example output: $2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG
	// Check the value. 
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == nil {
		fmt.Println("Passwords match")
	}
}
```
User model: 
```go
func (m *UserModel) Insert(name, email, password string) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPwd))
	if err != nil {
		var mySqlError *mysql.MySQLError
		// Using errors.As to check wether the error has the time *mysql.MySQLError. If so, assigning the error
		if errors.As(err, &mySqlError) {
			// If the error relates to our constraint, returning specific error
			if mySqlError.Number == 1062 && strings.Contains(mySqlError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}
```
