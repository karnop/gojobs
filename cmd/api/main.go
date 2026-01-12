package main 

import (
	"fmt"
	"log"
	"net/http"
)

// entry point of the application
func main() {
	// defining port here. in later iterations, will move it to env file
	const port = ":8080"

	// NewServeMux is a request multiplier (router).
	// It matches the URL of incoming request against a list of registered patterns
	// and calls the corresponding handler.
	mux := http.NewServeMux()

	// Register a simple health check route
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the GoJobs API")
	})

	log.Printf("Server started on port %s", port)

	// ListenAndServe starts an HTTP server with a given address and handler
	// This function blocks forever until the program is terminated
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}