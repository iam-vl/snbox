## Dependency injection

App struct must holds app-wide dependencies 
```go
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}
```
Update the handler definitions so they methods of app:
```go
func (app *application) HandleHome(w http.ResponseWriter, r *http.Request) {
    // code here
    // use the included properties
    // ... 
    app.errorLog.Print(err.Error)
    // ...
}
```
Update main:
```go
type application struct {}

func main() {
	port := flag.String("port", ":1111", "Server port")
	flag.Parse() // can use port as a flag

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}
	mux := http.NewServeMux()
    // Update necessary handlers
    // ... 
	mux.HandleFunc("/", app.HandleHome) // catch-all
	mux.HandleFunc("/snippet/view", app.HandleViewSnippet)
    // ... 
	// Custom http server
	srv := &http.Server{
		Addr:     *port,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Need to dereference a pointer
	infoLog.Printf("Starting server on port: %s", *port)
	// below mux ~ handler // mux is a special kind of handler
	err := srv.ListenAndServe()
	// err := http.ListenAndServe(*port, mux)
	errorLog.Fatal(err)
}
```

## Closures for dep injection 

If need to spread handlers across mult package, create package `config` exporting an `Application` struct, and have you handlers to form a closure. Roughly:

```go
func main() {
    app := &config.Application{
        ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    }
    mux.Handle("/", examplePackage.ExampleHandler(app))
    // ... 
}
func ExampleHandler(app *config.Application) http.HandlerFunc {
    return func(w http.RW, r *http.Req) {
        / ... 
        ts, err := template.ParseFiles(files...)
        if err @!= nil {
            app.ErrorLog.Print(err.Error())
            http.Error(w, "int server error", 500)
            return
        }
        // ...
    } 
}
```