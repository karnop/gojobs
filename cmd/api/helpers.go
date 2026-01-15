package main

import (
	"net/http"
)

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