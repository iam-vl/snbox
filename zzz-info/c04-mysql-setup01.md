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

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
```

```
curl -iL -X POST http://localhost:1111/snippet/create
```
DB.Exec() does the following:
1. Creates and comiples a prepared statement
2. Exec() passes the params to the statement
3. Deallocates the prepared statement on the db

## MySQL data conversions

```
CHAR, VARCHAR, TEXT -> string
BOOLEAN -> bool 
INT -> int; BIGINT -> int64
DECIMAL / NUMERIC -> float
TIME, DATE, TIMESTAMP -> time.Time
```

## Get by ID

```go
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := m.DB.QueryRow(query, id)
	s := &Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expired)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecords
		} else {
			return nil, err
		}
	}

	return s, nil
}
```

```
go run ./cmd/web
# github.com/iam-vl/snbox/internal/models
internal/models/snippets.go:45:16: undefined: ErrNoRecords
```

## Errors 

Recommended way: since 1.13 can add additional info by wrapping errors.
```go
if errors.Is(err, models.ErrNoRecord) {
	app.NotFound(w)
} else {
	app.ServerError(w, err)
}
```
Same as: 
```go
if err == models.ErrNoRecord {
	app.NotFound(w)
} else {
	app.ServerError(w, err)
}

```
Also: `errors.As()` - can check if an error has a specific type. 

## Manage null values

Scanning a null value: 

```go
err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
```
Res:
```
sql.Scan error ...
```

Solution: sqlNullString 
```go
type Book struct {
	Isbn string
	Title sql.NullString
	...
}
```

## Working with transactions 
Lock & unlock tables

```go
type ExampleModel struct {
	DB *sql.DB
}
func (m *ExampleModel) ExampleTransaction() error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // will roll back change if error
	query1 := `INSERT INTO...`
	query2 := `UPDATE...`

	_, err = tx.Exec(query1)
	if err != nil {
		return err
	}
	_, err = tx.Exec(query2)
	if err != nil {
		return err
	}
	// Commit if no errors
	err = tx.Commit() // Always call Rollback() or Commit() before yr function returns
	return err
}
```

## Prepared statements 

Create prepared statement using `DB.Prepare()`. Can embed in the model:
```go
type ExampleModel struct {
	DB *sql.DB
	InsertStmt *sql.Stmt
}
func NewExampleModel(db *sql.DB) (*ExampleModel, error) {
	insertStmt, err := db.Prepare("INSERT INTO...")
	if err != nil {
		return nil, err
	}
	return &ExampleModel{db, insertStmt}, nil
} 
func (m *ExampleModel) Insert(args...) error {
	// Calling Exec directly against the prepared statement
	// Works for Query and QueryRow too
	_, err := m.InsertrStmt.Exec(args...)
	return err
}
```
Using prep statement:
```go 
func main() {
	db, err := sql.Open(...)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	exampleModel, err := NewExampleModel(db) // mnew exampleModel includes the prepared statement
	if err != nil {
		errorLog.Fatal(err)
	}
	// -> prepared statement will be properly closed
	defer exampleModel.InsertStmt.Close() 
}
```
>**Note:**
>*Heavy load notice: Go uses a pool of many db connections.\
>When used the first tiome the pre stmt gets createwd on a db connection. 
>Sql.stmt remembers the connection. If it closes -> have to reprepare
>For the most part, simple Exec, Query, QueryRow are a good starting point