package data

import (
	"github.com/karnop/gojobs/internal/validator"
	"time"
	"fmt"
	"context"
	"database/sql"
)

// job represents a job posting in the application
type Job struct {
	// omitempty means if id is empty, hide it in JSON
	Id          int `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Company     string `json:"company"`
	Salary      int    `json:"salary"`
	UserId      int    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type JobModel struct {
	DB *sql.DB
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
	v.Check(job.Salary >= 0, "salary", "must be a positive number")
}


// Insert adds a new job to the database
func (m JobModel) Insert(job *Job) error {
	query := `
		INSERT INTO jobs (title, description, company, salary, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at`

	// Use QueryRow because we want to get the ID back
	return m.DB.QueryRow(query, job.Title, job.Description, job.Company, job.Salary, job.UserId).Scan(&job.Id, &job.CreatedAt)
}

// Get fetches a single job by ID
func (m JobModel) Get(id int) (*Job, error) {
	query := `
		SELECT id, title, description, company, salary, user_id, created_at
		FROM jobs
		WHERE id = $1`

	var job Job
	err := m.DB.QueryRow(query, id).Scan(
		&job.Id,
		&job.Title,
		&job.Description,
		&job.Company,
		&job.Salary,
		&job.UserId,
		&job.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}

// GetAll fetches a list of jobs based on filters
func (m JobModel) GetAll(title string, company string, filters Filters) ([]*Job, error) {
	query := fmt.Sprintf(`
		SELECT id, title, company, description, salary, user_id, created_at
		FROM jobs
		WHERE (LOWER(title) LIKE LOWER($1) OR $1 = '')
		AND (LOWER(company) LIKE LOWER($2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	// prepare arguments
	args := []interface{}{
		"%" + title + "%",     // $1: search title (partial match)
		"%" + company + "%",   // $2: search company (partial match)
		filters.limit(),       // $3: limit
		filters.offset(),      // $4: offset
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// executing Query
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	jobs := []*Job{}
	for rows.Next() {
		var job Job
		err := rows.Scan(
			&job.Id,
			&job.Title,
			&job.Company,
			&job.Description,
			&job.Salary,
			&job.UserId,
			&job.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}