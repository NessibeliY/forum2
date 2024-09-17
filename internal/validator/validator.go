package validator

import (
	"regexp"
	"unicode"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid returns true if the errors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map (so long as no entry already exists for the given key).
func (v *Validator) AddErrors(key, message string) {
	if _, exist := v.Errors[key]; !exist {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddErrors(key, message)
	}
}

// Validated email
func (v *Validator) Matches(value string, regx *regexp.Regexp) bool {
	return regx.MatchString(value)
}

// Check min len password
func (v *Validator) MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func (v *Validator) ValidPassword(password string) bool {
	hasDigit := false
	hasUpper := false
	hasLower := false

	for _, ch := range password {
		if unicode.IsDigit(ch) {
			hasDigit = true
		} else if unicode.IsUpper(ch) {
			hasUpper = true
		} else if unicode.IsLower(ch) {
			hasLower = true
		}
	}

	if !hasDigit || !hasUpper || !hasLower {
		return false
	}

	return true
}
