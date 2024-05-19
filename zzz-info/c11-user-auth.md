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

## User login 

Update validator to support valid errors that are not associated with one specific field. Later,  can show something like "your email or pwd is incorrect". 
Add NonFieldError to validator.go
```go
type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
}

func (v *Validator) Valid8() bool {
	// Include NonFieldErrors in the validation
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// Add a NonFieldError to the NonFieldError slice
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}
```
Add login template (login.tmpl). 
In handlers, create `UserLoginForm` struct and update `HandleLoginForm()`:
```go
type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
func (app *application) HandleLoginForm(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = UserLoginForm{}
	app.Render(w, http.StatusOK, "login.tmpl", data)
}
```
## Verifying user details 

UserModel.Auth():
```go
func (m *UserModel) Auth(email, password string) (int, error) {
	var id int
	var pwdHash []byte
	stmt := `SELECT id, hashed_pwd FROM users WHERE email = ?`
	// check for creds
	err := m.DB.QueryRow(stmt, email).Scan(&id, &pwdHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCreds // Not found
		} else {
			return 0, err
		}
	}
	// if found
	err = bcrypt.CompareHashAndPassword(pwdHash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCreds // incorrect pwd
		} else {
			return 0, err
		}
	}
	return id, nil
}
```
Handler:  
```go
func (app *application) HandleLoginPost(w http.ResponseWriter, r *http.Request) {
	var form UserLoginForm
	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be a valid address")
	if !form.Valid8() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
	}
	// check if the creds are valid
	id, err := app.users.Auth(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCreds) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.ServerError(w, err)
		}
		return
	}
	// Generate a new session ID when the auth status and priovilege level change
	// For example, if user login / logout.
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// Add the ID of current user to session, so they are now logged in.
	app.sessionManager.Put(r.Context(), "authenticatedUseId", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther) // 403
}
```

## Logging out

Essentially, remove `authenticatedUserID` value from the session + redirect: 
```go
func (app *application) HandleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// remove the user id from the session data
	app.sessionManager.Remove(r.Context(), "authenticatedUserId")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

## User authorization 

Goals:  
1. Only authenticated (i.e. logged in) users can create a new snippet; and
2. The contents of the navigation bar changes depending on whether a user is authenticated
(logged in) or not. Specifically:
  * Authenticated users should see links to ‘Home’, ‘Create snippet’ and ‘Logout’.
  * Unauthenticated users should see links to ‘Home’, ‘Signup’ and ‘Login’.
Solution:   
1. Check if the reqt comes from an authd user: create helper: `app.isAuthenticated(r *http.Request) bool {}`. 
2. Pass the info to the HTML templates:  
  2.1. Add `isAuth` to template data:
  2.2. Update `NewTemplateData()` helper so the if is automcatically added to templateData when we render a template. 

Helper:  
```go
func (app *application) IsAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserId")
}
```
Tempaltes:
```go
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string // Flash message
	isAuth      bool   // Add to templ data
}
```
Helpers:
```go
func (app *application) NewTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(r.Context(), "flash"),
		// Add the auth status to the templ data
		IsAuth: app.IsAuthenticated(r),
	}
}
```  
Add to the templates (nav):
```go
{{ define "nav" }}
<nav>
    <div>
        <a href="/">Home</a>
        {{if .IsAuth}}
            <a href="/snippet/create">Create snippet</a>
        {{end}}
        
    </div>
    <div>
        {{if .IsAuth}}
            <form action="/user/logout" method="POST">
                <button>Logout</button>
            </form>
        {{else}}
            <a href="/user/signup">Signup</a>
            <a href="/user/login">Log in</a>
        {{end}}
    </div>    
</nav>
{{ end }}
```
## Restrict access to the snippet form 

Create middleware:
```go
func (app *application) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// Otherwise, set "Cache-Control: no store" header
		// so that the pages that require auth are not stored in the user browser cache
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
```
Add the new mw to routes. Rearrange: 
* unprotected routes use existing dynamic chain
* protected routes: existing chain + new middleware 
```go
func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})

	// static file server
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))

	// Middleware chain that will contain
	dynamic := alice.New(app.sessionManager.LoadAndSave)
	// Protected middleware chain:
	protectedChain := dynamic.Append(app.RequireAuth)

	// Unprotected routes
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.HandleHome)) // catch-all
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.HandleViewSnippet))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.HandleSignupForm))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.HandleSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.HandleLoginForm))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.HandleLoginPost))
	// Protected routes 
	router.Handler(http.MethodGet, "/snippet/create", protectedChain.ThenFunc(app.HandleSnippetForm))
	router.Handler(http.MethodPost, "/snippet/create", protectedChain.ThenFunc(app.HandleCreateSnippet))	
	router.Handler(http.MethodPost, "/user/logout", protectedChain.ThenFunc(app.HandleLogoutUser))

	mwareChain := alice.New(app.RecoverPanic, app.LogRequest, SecureHeaders)
	return mwareChain.Then(router)
}
```
Example w/out alice:
```go
router.Handler(http.MethodPost,"/snippet/create",app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.snippetCreate)))
```

Can check thru `curl`:
```sh
$ curl -ki -X POST https://localhost:1111/snippet/create
HTTP/2 303 
content-security-policy: default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com
location: /user/login
referrer-policy: origin-when-cross-origin
vary: Cookie
x-content-type: nosniff
x-frame-options: deny
x-xss-protection: 0
content-length: 0
date: Sun, 19 May 2024 10:47:21 GMT
```


