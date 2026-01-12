package main 

import (
	"encoding/json"
	"net/http"
	"github.com/karnop/gojobs/internal/data"
)


// createJobHandler handles POST request to add a new job
func (app *application) createJobHandler(w http.ResponseWriter, r *http.Request) {
	// variable to hold the incoming data
	var job data.Job

	// decoding json body from request
	err := json.NewDecoder(r.Body).Decode(&job) 
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return 
	}

	// SQL query (postgres)
	// RETURNING id allows us to get the auto-generated ID back from SQLSERVER.
	// $1, $2, $3, $4 are placeholders for the data to prevent SQL Injection.
	query := `
		INSERT INTO jobs (title, description, company, salary)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	// QueryRow executes a query that returns exactly one row (the ID).
	err = app.DB.QueryRow(query, job.Title, job.Description, job.Company, job.Salary).Scan(&job.Id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// responding to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}


// listjobshandler handles GET request to show all jobs
func (app *application) listJobsHandler(w http.ResponseWriter, r *http.Request) {
	// simple select
	query := "SELECT id, title, description, company, salary FROM jobs"
	rows, err := app.DB.Query(query)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}
	// closing rows to free up db connection 
	defer rows.Close()

	var jobs []data.Job
	for rows.Next() {
		var j data.Job
		// Scan copies the columns from the current row into the values pointed at.
		err := rows.Scan(&j.Id, &j.Title, &j.Description, &j.Company, &j.Salary)
		if err != nil {
			http.Error(w, "Error scanning database row", http.StatusInternalServerError)
			return
		}
		jobs = append(jobs, j)
	}

	// errors that might have occurred during iteration
	if err = rows.Err(); err != nil {
		http.Error(w, "Database iteration error", http.StatusInternalServerError)
		return
	}	

	// setting the header
	w.Header().Set("Content-Type", "application/json")

	// encoding jobs slice directly to the response writer
	err = json.NewEncoder(w).Encode(jobs)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}