package main 

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"
	"github.com/joho/godotenv"
	"os"
)


// defining a struct to hold application dependencies 
// it makes the handlers cleaner because they can access the DB via this struct
type application struct {
	DB *sql.DB
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
		DB: db,
	}

	// NewServeMux is a request multiplier (router).
	// It matches the URL of incoming request against a list of registered patterns
	// and calls the corresponding handler.
	mux := http.NewServeMux()

	// Register a simple health check route
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the GoJobs API")
	})

	// GET /jobs 
	mux.HandleFunc("GET /jobs", app.listJobsHandler)

	// POST /jobs
	mux.HandleFunc("POST /jobs", app.createJobHandler)

	// GET /job/id
	mux.HandleFunc("GET /jobs/{id}", app.getJobHandler)

	log.Printf("Server started on port %s", port)

	// ListenAndServe starts an HTTP server with a given address and handler
	// This function blocks forever until the program is terminated
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
