package postgres

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

// InitConnection initializes and returns a connection to a PostgreSQL database.
//
// It uses the POSTGRES_DSN environment variable to establish the connection.
// The function attempts to open a connection and verify it with a ping.
//
// Returns:
//   - *sql.DB: A pointer to the database connection if successful.
//   - error: An error if the connection fails to open or ping.
func InitConnection() (*sql.DB, error) {
    db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
