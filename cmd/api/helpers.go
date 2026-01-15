package main

import (
	"net/http"
	"net/url"
    "strconv"
	"github.com/karnop/gojobs/internal/validator"
)

// readString returns a string value from the query string, or the default value
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

// readInt returns an integer value from the query string, or the default value
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

// serverError logs the detailed error and sends a generic 500 to the user
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	// We include the request method and URL so we know WHERE it happened.
	app.Logger.Error("server error", 
		"method", r.Method, 
		"url", r.URL.String(), 
		"error", err.Error(),
	)

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// clientError sends a specific status code and description to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}