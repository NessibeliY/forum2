package validator

import (
	"regexp"
	"unicode"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("(?:[a-z0-9!#$%&'*+\\/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+\\/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")

type Validator struct {
	ErrorsMap map[string]string
}

func NewValidator() *Validator {
	return &Validator{
		ErrorsMap: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.ErrorsMap) == 0
}

// AddError adds an error message to the map (as long as no entry already exists for the given key).
func (v *Validator) AddErrors(key, message string) {
	if _, exist := v.ErrorsMap[key]; !exist {
		v.ErrorsMap[key] = message
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
	hasDigit, hasUpper, hasLower := false, false, false

	for _, ch := range password {
		switch {
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		}

		if hasDigit && hasUpper && hasLower {
			return true
		}
	}

	return hasDigit && hasUpper && hasLower
}
