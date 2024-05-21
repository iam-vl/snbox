# CH 13 Optional Go Features 

Plan: 
1. File embedding
2. Generics 

## Embedding 

Create file (ui/efs.go):
```go
package ui
import "embed"
//go:embed "html" "static"
var Files embed.FS
```
About the comment directive: 
* During compilation, instructs Go to store the files from ui/html and ui/static in an `embed.FS` embedded file system.
* Must be placed immediately above the variabl 
* Supports multiple paths
* Can only use the directive onb global variables at package level, not within functions or methods. 
* If used inside a func => error `go:embed cannot apply to var inside 
About paths: 
Cannot contain . or .. elements and may not begin or end with /. Basically, restricts you to only embedding files in the same directory as the source code that contains the directive. 
If path to dir, all files recursively emnbedded, excepts the opnes that begin with '.' or '_'. To include those files: `go:embed "all:static`. 
Path separator only forward slash
Embedded fs always rooted to the dir that contains the directive. 

## Updating the app: serve the statics from FS

Routes (update the file server): 
```go
func main() {
    // ..

    // Convert ui.Files embedded fs to a http.FS type so it works as a http.FileSystem interface
	// Then pass it to http.FileServer to create a file (server) handler
	fileServer := http.FileServer(http.FS(ui.Files))
	// fileserver := http.FileServer(http.Dir("./ui/static/"))
	// Our statics are now in the static of folder of the embedded fs. We no longer need to strip the prefix.
	// Any requests with `/static/` will now be passed directly to file server.
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	// router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))
    
}
```

