package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // Not using it, but need the init() function
	"github.com/iam-vl/snbox/internal/models"
)

const (
	pwd = "vl#123pass"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	port := flag.String("port", ":1111", "Server port")
	dsnText := fmt.Sprintf("web:%s@/snbox?parseTime=true&allowNativePasswords=true", pwd)
	dsn := flag.String("dsn", dsnText, "sb_mysql_datasource")
	flag.Parse() // can use port as a flag

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDb(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := NewTemplateCache()
	// templateCache, err := NTCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Custom http server
	srv := &http.Server{
		Addr:     *port,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Need to dereference a pointer
	infoLog.Printf("Starting server on port: %s", *port)
	// below mux ~ handler // mux is a special kind of handler
	err = srv.ListenAndServe()
	// err := http.ListenAndServe(*port, mux)
	errorLog.Fatal(err)
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn) // Initializing connection pool
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Ping successful")
	fmt.Printf("Database: %+v\n", db)
	return db, nil
}
