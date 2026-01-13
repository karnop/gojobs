package data

import (
	"github.com/karnop/gojobs/internal/validator"
)

// job represents a job posting in the application
type Job struct {
	// omitempty means if id is empty, hide it in JSON
	Id string `json:"id,omitempty"` 
	Title string `json:"title"`
	Description string `json:"description"`
	Company string `json:"company"`
	Salary string `json:"salary"`
}

// ValidateJob checks if the Job struct is safe to insert
func ValidateJob(v *validator.Validator, job *Job) {
	// title must not be empty and less than 100 chars
	v.Check(job.Title != "", "title", "must be provided")
	v.Check(len(job.Title) <= 100, "title", "must not be more than 100 characters")

	// Description is required
	v.Check(job.Description != "", "description", "must be provided")

	// Company is required
	v.Check(job.Company != "", "company", "must be provided")

	// Salary is required
	v.Check(job.Salary != "", "salary", "must be provided")
}