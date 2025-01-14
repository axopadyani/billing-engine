package postgres

import (
	"context"
	"database/sql"

	"go.uber.org/multierr"
)

// executor is an interface for database executor, which should be implemented by *sql.DB and *sql.Tx.
type executor interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// finishTransaction completes a database transaction based on the presence of an error.
// If an error is present, it rolls back the transaction. Otherwise, it commits the transaction.
//
// Parameters:
//   - err: The error to check. If not nil, the transaction will be rolled back.
//   - tx: The SQL transaction to be committed or rolled back.
//
// Returns:
//   - An error if the rollback or commit fails, or the original error if rollback succeeds, otherwise nil.
func finishTransaction(err error, tx *sql.Tx) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return multierr.Combine(err, rollbackErr)
		}

		return err
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			return commitErr
		}

		return nil
	}
}
