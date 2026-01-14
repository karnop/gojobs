package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/karnop/gojobs/internal/data"
	"log"
	"net/http"
	"os"
)

// defining a struct to hold application dependencies
// it makes the handlers cleaner because they can access the DB via this struct
type application struct {
	DB    *sql.DB
	Users data.UserModel
}

// entry point of the application
func main() {
	// loading the env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	port := os.Getenv("port")
	if port == "" {
		log.Fatal("DB_DSN environment variable not set")
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable not set")
	}
	log.Println("Connecting to Cloud Database...")

	// calling helper function to open the connection
	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer db.Close()
	log.Println("SUCCESS! Database connection established")

	app := &application{
		DB:    db,
		Users: data.UserModel{DB: db},
	}

	// NewServeMux is a request multiplier (router).
	// It matches the URL of incoming request against a list of registered patterns
	// and calls the corresponding handler.
	mux := http.NewServeMux()

	// Register a simple health check route
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the GoJobs API")
	})

	mux.HandleFunc("GET /jobs", app.listJobsHandler)
	mux.HandleFunc("POST /jobs", app.authenticate(app.createJobHandler))
	mux.HandleFunc("GET /jobs/{id}", app.getJobHandler)
	mux.HandleFunc("POST /users", app.registerUserHandler)
	mux.HandleFunc("POST /users/login", app.loginUserHandler)

	log.Printf("Server started on port %s", port)

	// ListenAndServe starts an HTTP server with a given address and handler
	// This function blocks forever until the program is terminated
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
