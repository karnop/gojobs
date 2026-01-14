package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// openDB opens a connection to the database
// param dsn : Data Source Name (The address string of the DB)
// it returns a generic *sql.DB pool (The handle to use the DB) that is safe for concurrent use
func openDB(dsn string) (*sql.DB, error) {
	// open the connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// testing the connection
	// sql.open() doesnt connect immediately, it is lazy
	// db.Ping() forces a connection to ensure DB is actually alive
	// If the DB is down or password is wrong, this will fail.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
