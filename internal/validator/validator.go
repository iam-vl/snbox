package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func MinChars(val string, n int) bool {
	// True if a val contains at least n chars
	return utf8.RuneCountInString(val) >= n
}
func Matches(val string, rx *regexp.Regexp) bool {
	// True if the val matches the regex
	return rx.MatchString(val)
}

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
