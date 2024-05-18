package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // Not using it, but need the init() function
	"github.com/iam-vl/snbox/internal/models"
)

const (
	pwd = "vl#123pass"
)

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

	// templateCache, err := NewTemplateCache()
	templateCache, err := NewTemplateCache3()
	// templateCache, err := NTC2() // before
	// templateCache, err := NTC()
	if err != nil {
		fmt.Println("Template cache error!!")
		errorLog.Fatal(err)
	}

	// Initialize a decoder
	formDecoder := form.NewDecoder()
	// Configure a sesh manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Serve over https
	sessionManager.Cookie.Secure = true

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
		Addr: *port,
		// MaxHeaderBytes: 524288, // Limit header size. If above: 431 Request Header Fields Too Large
		ErrorLog: errorLog,
		Handler:  app.routes(),
		// Add custom TLS config
		TLSConfig: tlsConfig,
		// Add Idle, Read, and Write timeouts
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Need to dereference a pointer
	infoLog.Printf("Starting server on port: %s", *port)
	// Use TLS
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	// err = srv.ListenAndServe()
	// err := http.ListenAndServe(*port, mux) // legacy
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
