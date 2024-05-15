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

## Optimize HTTPS 

Restrict some elliptic curves. to do so, create a `tls.Config` struct with custom TLS settings and add it to `http.Server`. Main: 
```go
func main() {

	// Set up everything
    // ... 
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
	// create somewhere to hold custom TLS settings 
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
	}

	// Custom http server
	srv := &http.Server{
		Addr:     *port,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		// Add custom TLS config
        TLSConfig: tlsConfig,
	}

	// Need to dereference a pointer
	infoLog.Printf("Starting server on port: %s", *port)
	// Use TLS
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	// err = srv.ListenAndServe()
	// err := http.ListenAndServe(*port, mux) // legacy
	errorLog.Fatal(err)
}

To config min and max TLS versions:
```go 
tlsConfig := &tls.Config {
    MinVersion: tls.VersionTLS12,
    MaxVersion: tls.VersionTLS13, // so far the latest one. 
}
```
Can restrict cipher suites. For example, only support cipher suites that use ECDHE (forward security) and not support weak suites:
```go
tlsConfig := &tls.Config {
    CipherSuites: []uint16{
        tls.TLS_ECDHE_ECDSA_WITH_AES256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES256_GCM_SHA384,
        tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
        tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
        tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
    }
}
```

>**Note**
>If a TLS 1.3 connection is negotiated, any CipherSuites field in your tls.Config will be ignored. The reason for this is that all the cipher suites that Go supports for TLS 1.3 connections are considered to be safe, so there isn’t much point in providing a mechanism to configure them.

## Connection timeouts 

Improve server resiliency by adding timeout settings:  
```go
	srv := &http.Server{
		Addr: *port,
		// MaxHeaderBytes: 524288, // Limit header size. If above: 431 Request Header Fields Too Large (=4996 bytes)
		ErrorLog: errorLog,
		Handler:  app.routes(),
		// Add custom TLS config
		TLSConfig: tlsConfig,
		// Add Idle, Read, and Write timeouts
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
```
All three are server-wide settings which act on underlying connection and apply to all requests. 
Go enables keep-alives (persistent conn). Automatically closed after a couple minutes. You cannot increase,  but can decrease the wait time. 

>**Note**
>If you set ReadTimeout, but don't set IdleTimeout, Idle will default to ReadTimeout. 
>Prevent the data that handler reaturns from taking too long to write with WriteTimeout. 

>**Note**
>
