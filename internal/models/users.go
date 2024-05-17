package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int
	Name       string
	Email      string
	HanshedPwd []byte
	Created    time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, pwd string) error {
	return nil
}
func (m *UserModel) Auth(email, pwd string) (int, error) {
	return 0, nil
}
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
