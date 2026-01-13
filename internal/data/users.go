package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/karnop/gojobs/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgx/v5/pgconn" // to handle Postgres specific errors
)

// custom error 
var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// User represents a registered user
type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"` // - means never send in JSON
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// password is a custom struct to handle hashing logic
type password struct {
	plaintext *string
	hash []byte
}

// set calculates the bcrypt hash of a plaintext password
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

// Matches checks if a plaintext password matches the stored hash
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// UserModel wraps the DB connection pool
type UserModel struct {
	DB *sql.DB
}

// insert adds a new user to the db
func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	args := []interface{}{user.Name, user.Email, user.Password.hash, user.Role}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt)
	if err != nil {
		// We use errors.As to check if the error is a specific Postgres error type
		var pgErr *pgconn.PgError 
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Code "23505" is the official SQL State for "unique_violation"
				return ErrDuplicateEmail
			}
		}

		return err
	}
	return nil
}

// ValidateUser checks the request data.
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(user.Email != "", "email", "must be provided")
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be a valid email address")

	if user.Password.plaintext != nil {
		v.Check(*user.Password.plaintext != "", "password", "must be provided")
		v.Check(len(*user.Password.plaintext) >= 8, "password", "must be at least 8 bytes long")
	}
}