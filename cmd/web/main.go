package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", ":1111", "Server port")
	flag.Parse() // can use port as a flag
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	mux := http.NewServeMux()
	// os_port := os.Getenv("SBOX_PORT")

	// static file server
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("/book", HandleDownloader)

	mux.HandleFunc("/", HandleHome) // catch-all
	mux.HandleFunc("/snippet/view", HandleViewSnippet)
	mux.HandleFunc("/snippet/create", HandleCreateSnippet)
	mux.HandleFunc("/head", HandleCustomizeHeaders)
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
