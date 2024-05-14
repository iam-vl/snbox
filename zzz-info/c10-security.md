# CH 10: Security improvements 

## Plan 

1. Create self-signed TLS certificate.
2. Serve the app over HTTPS.
3. Tweak default TLS settings. 
4. Set conn timeouts to mitigate slow-client attacks. 

## Generating TLS certs 

Go lib includes the `generate_cert.go` tool. 

```sh
mkdir tls; cd tls
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

## Turn on TLS
Main: 
```go
func main() {

	port := flag.String("port", ":1111", "Server port")
	dsnText := fmt.Sprintf("web:%s@/snbox?parseTime=true&allowNativePasswords=true", pwd)
	dsn := flag.String("dsn", dsnText, "sb_mysql_datasource")
	flag.Parse() // can use port as a flag
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, _ := openDb(*dsn)
	defer db.Close()

	templateCache, _ := NTC2()
    // 
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
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
	srv := &http.Server{
		Addr:     *port,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on port: %s", *port)
	// Use TLS
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
```
Firefox: **Accept risks and continue**. 
See sec data: Ctrl + I
Http/2 (if supported by client) => Will load faster

Gencert tool grants: 
* Read permission for all users for cert.pem
* Read permission only to the owner of key.pem 

```sh
$ ls -la tls
-rw-rw-r-- 1 dell dell 1090 May 14 22:31 cert.pem
-rw------- 1 dell dell 1704 May 14 22:31 key.pem
```

