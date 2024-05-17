package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching records found")

	// Use this err later if a user tries to log in w/ an incorrect email address / pwd
	ErrInvalidCreds = errors.New("models: invalid creds")
	// Existing email error
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
