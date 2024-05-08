package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid8() bool {
	return len(v.FieldErrors) == 0
}

// Add an err message to FE map (if no entry for this key)
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// Add an error msg to the map only if a validation check not okay
func (v *Validator) CheckField(ok bool, key, msg string) {
	if !ok {
		v.AddFieldError(key, msg)
	}
}

func NotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}

func MaxChars(s string, n int) bool {
	return utf8.RuneCountInString(s) <= n
}

// Returns true if in a list of permitted ints
func PermittedInt(val int, permittedVals ...int) bool {
	for i := range permittedVals {
		if val == permittedVals[i] {
			return true
		}
	}
	return false
}
