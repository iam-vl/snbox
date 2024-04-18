# Setting up mysql

See `c0401-setup-db.sql`.

Can upgrade the version

```
go get github.com/go-sql-driver/mysql
go get github.com/go-sql-driver/mysql@v1 // latest v1
```

## Snippet model 

```go
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expired time.Time
}
type SnippetModel struct {
	DB *sql.DB
}
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}
```

## DB query methods 

```go
DB.Query()  // SELECT, multiple rows
DB.QueryRow() // SELECT, single row
DB.Exec() // No returns 
```

## Create new

```go
func (app *application) HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		app.ClientError(w, http.StatusMethodNotAllowed) // using ClientError()
		return
	}
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
```

```
curl -iL -X POST http://localhost:1111/snippet/create
```

