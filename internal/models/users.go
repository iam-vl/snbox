package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	fmt.Println("inserting...")
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_pwd, created) VALUES (?, ?, ?, UTC_TIMESTAMP())`
	fmt.Printf("Creds (inc pwd hash): %s, %s, %s\n", name, email, pwdHash)
	_, err = m.DB.Exec(stmt, name, email, string(pwdHash))
	fmt.Println("Insert user model 2")
	if err != nil {
		fmt.Println("Insert user model 3")
		var mySqlError *mysql.MySQLError
		// Using errors.As to check wether the error has the time *mysql.MySQLError. If so, assigning the error
		if errors.As(err, &mySqlError) {
			// If the error relates to our constraint, returning specific error
			if mySqlError.Number == 1062 {
				// if mySqlError.Number == 1062 && strings.Contains(mySqlError.Message, "users_uc_email") {
				fmt.Println("Yes, ErrDuplicateEmail")
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}
func (m *UserModel) Auth(email, password string) (int, error) {
	var id int
	var pwdHash []byte
	stmt := `SELECT id, hashed_pwd FROM users WHERE email = ?`
	// check for creds
	err := m.DB.QueryRow(stmt, email).Scan(&id, &pwdHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCreds
		} else {
			return 0, err
		}
	}
	// if found
	err = bcrypt.CompareHashAndPassword(pwdHash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCreds
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
