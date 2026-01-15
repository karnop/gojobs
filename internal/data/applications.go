package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// ErrDuplicateApplication is returned when a user applies twice
var ErrDuplicateApplication = errors.New("you have already applied to this job")

type JobApplication struct {
	Id int `json:"id"`
	JobId int `json:"job_id"`
	UserId int `json:"user_id"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type JobApplicationModel struct {
	DB *sql.DB
}

// Insert creates a new application record
func(m JobApplicationModel) Insert(application *JobApplication) error {
	query := `
		INSERT INTO applications (job_id, user_id, status)
		VALUES ($1, $2, 'applied')
		RETURNING id, created_at, status
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, application.JobId, application.UserId).Scan(
		&application.Id,
		&application.CreatedAt,
		&application.Status,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Check for Unique Violation (Code 23505)
			if pgErr.Code == "23505" {
				return ErrDuplicateApplication
			}
		}
		return err
	}

	return nil
}