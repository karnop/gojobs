package main 

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/karnop/gojobs/internal/data"
)

// jobs is our temporary in memory db
var jobs = []data.Job{
	{
		Id:          "1",
		Title:       "Software Engineer",
		Description: "Go developer needed",
		Company:     "Google",
		Salary:      "$120,000",
	},
} 

// createJobHandler handles POST request to add a new job
func createJobHandler(w http.ResponseWriter, r *http.Request) {
	// variable to hold the incoming data
	var job data.Job

	// decoding json body from request
	err := json.NewDecoder(r.Body).Decode(&job) 
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return 
	}

	// assigning a unique id using time
	job.Id = time.Now().Format("20060102150405")

	// adding the new job to db
	jobs = append(jobs, job)

	// responding to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}


// listjobshandler handles GET request to show all jobs
func listJobsHandler(w http.ResponseWriter, r *http.Request) {
	// setting the header
	w.Header().Set("Content-Type", "application/json")

	// encoding jobs slice directly to the response writer
	err := json.NewEncoder(w).Encode(jobs)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}