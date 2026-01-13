package validator

import (
	"regexp"
)

// EmailRx is the standard regex for validating email formats
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")


// Validator contains a map of validation errors
// the key is the field name, and the value is the error message
type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator {
		Errors: make(map[string]string),
	}
}

// Valid returns true if the Errors map is empty
func (v * Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map if it doesnt exist
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check is a helper
// if ok is false, it add an error

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// returns true if a string value matches a specific regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}