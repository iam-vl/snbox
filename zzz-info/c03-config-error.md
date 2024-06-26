## CH 03: Config & error handling 

Learn to: 
* Set configs at runtime using command-line flags
* Improve app log messages  
* Make dependencies available to handlers in an extensible typesafe way 
* Centralize error handling 

## Pre-existing variables 

```go
type config struct {
    port string
    staticDir string
}
...
var cfg config 
flag.StringVar(&cfg.port, "port", ":1234", "TCP server port")
flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets rel project root")
flag.Parse()
```

## Custom loggers 

```go
...
infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
...
infoLog.Printf("Starting server on port: %s", *port)
err := http.ListenAndServe(*port, mux)
errorLog.Fatal(err)
```
Can also use `log.Llongtime` and `log.LUTC flags`. 

## Logging to a file 

```sh
go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log
```
```go
f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
if err != nil {
    log.Fatal(err)
}
defer f.Close()
infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
```

## Custom HTTP server 

```go
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
```
