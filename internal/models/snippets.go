package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

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
	query := `INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// Use the Exec method on the embedded connection pool to execute the statement
	// Perfectly fine to ignore the res, if you don't need it
	result, err := m.DB.Exec(query, title, content, expires)
	fmt.Printf("Result type: %T\n", result)
	fmt.Printf("Result val: %+v\n", result)
	if err != nil {
		return 0, err
	}
	// get the new record's id
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
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

func (m *SnippetModel) Latest10() ([]*Snippet, error) {
	return nil, nil
}
