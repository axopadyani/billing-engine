package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/shopspring/decimal"

	"github.com/axopadyani/billing-engine/internal/entity"
)

// Repository represents a data access layer for interacting with a PostgreSQL database.
// It encapsulates database operations and provides methods for querying and manipulating data.
type Repository struct {
	// db is a pointer to the SQL database connection.
	// It is used to execute SQL queries and transactions.
	db *sql.DB
}

// NewRepository creates and returns a new Repository instance.
//
// It takes a database connection as an argument and initializes the Repository
// with this connection.
//
// Parameters:
//   - db: A pointer to sql.DB representing the database connection to be used by the repository.
//
// Returns:
//   - A pointer to a new Repository instance initialized with the provided database connection.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateLoan creates a new loan record in the database.
// It performs the operation within a transaction to ensure data consistency.
//
// The function executes the following steps:
// 1. Starts a new transaction with serializable isolation level.
// 2. Retrieves the latest loan for the user.
// 3. Validates the new loan using the provided validation function.
// 4. Inserts the new loan into the database if validation passes.
//
// Parameters:
//   - ctx: A context.Context for handling cancellation and timeouts.
//   - loan: An entity.Loan instance containing the loan details to be created.
//   - validateFn: A function that takes the latest loan as an argument and returns an error if validation fails.
//
// Returns:
//
//	An error if any step in the process fails, including database errors, validation errors,
//	or transaction errors. Returns nil if the loan is successfully created.
func (r *Repository) CreateLoan(
	ctx context.Context,
	loan *entity.Loan,
	validateFn func(latestLoan *entity.Loan) error,
) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer func() { err = finishTransaction(err, tx) }()

	latestLoan, err := getLatestLoan(ctx, tx, loan.UserID)
	if err != nil {
		return err
	}

	if err := validateFn(latestLoan); err != nil {
		return err
	}

	query, args := loanStruct.InsertInto(loansTable, toPostgresLoan(loan)).BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

// GetLatestLoan retrieves the most recent loan for a given user from the database.
//
// This function constructs and executes a SQL query to fetch the latest loan
// based on the creation timestamp for the specified user ID. It returns the loan
// as an entity.Loan object or nil if no loan is found.
//
// Parameters:
//   - ctx: A context.Context for handling cancellation and timeouts.
//   - userID: The UUID of the user whose latest loan is being retrieved.
//
// Returns:
//   - *entity.Loan: The most recent loan entity if found, or nil if no loan exists.
//   - error: An error object if any database operation fails, or nil if successful.
func (r *Repository) GetLatestLoan(ctx context.Context, userID uuid.UUID) (*entity.Loan, error) {
	return getLatestLoan(ctx, r.db, userID)
}

func getLatestLoan(ctx context.Context, executor executor, userID uuid.UUID) (*entity.Loan, error) {
	sb := loanStruct.SelectFrom(loansTable)
	query, args := sb.Where(sb.Equal("user_id", userID)).
		OrderBy("created_at").Desc().
		Limit(1).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	var pgLoan postgresLoan
	err := executor.QueryRowContext(ctx, query, args...).Scan(loanStruct.Addr(&pgLoan)...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return pgLoan.toEntityLoan(), nil
}

// GetLoanPaidAmount retrieves the total amount paid for a specific loan.
//
// This function constructs and executes a SQL query to calculate the sum of all
// payments made for the given loan ID.
//
// Parameters:
//   - ctx: A context.Context for handling cancellation and timeouts.
//   - loanID: The UUID of the loan for which to calculate the total paid amount.
//
// Returns:
//   - decimal.Decimal: The total amount paid for the loan. Returns decimal.Zero if no payments are found.
//   - error: An error object if any database operation fails, or nil if successful.
func (r *Repository) GetLoanPaidAmount(ctx context.Context, loanID uuid.UUID) (decimal.Decimal, error) {
	return getLoanPaidAmount(ctx, r.db, loanID)
}

func getLoanPaidAmount(ctx context.Context, executor executor, loanID uuid.UUID) (decimal.Decimal, error) {
	sb := sqlbuilder.NewSelectBuilder()
	query, args := sb.Select("SUM(amount)").
		From(loanPaymentsTable).
		Where(sb.Equal("loan_id", loanID)).
		GroupBy("loan_id").
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	var sum decimal.Decimal
	err := executor.QueryRowContext(ctx, query, args...).Scan(&sum)
	if errors.Is(err, sql.ErrNoRows) {
		return decimal.Zero, nil
	} else if err != nil {
		return decimal.Zero, err
	}

	return sum, nil
}

// MakePayment processes a payment for a loan, updates the loan record if necessary, and returns the updated loan information.
//
// This function performs the following operations within a transaction:
// 1. Retrieves the loan information.
// 2. Calculates the current paid amount for the loan.
// 3. Executes the provided makePaymentFn to process the payment.
// 4. Inserts a new loan payment record.
// 5. Updates the loan record if required.
//
// Parameters:
//   - ctx: A context.Context for handling cancellation and timeouts.
//   - loanID: The UUID of the loan for which the payment is being made.
//   - paymentAmount: The amount of the payment being made, as a decimal.Decimal.
//   - makePaymentFn: A function that processes the payment, determines if the loan should be updated,
//     and returns the payment details. It takes the current loan and paid amount as arguments.
//
// Returns:
//   - loan: An entity.Loan instance representing the updated loan information.
//   - newPaidAmount: A decimal.Decimal representing the new total paid amount for the loan after this payment.
//   - err: An error object if any step in the process fails, or nil if the payment is successfully processed.
func (r *Repository) MakePayment(
	ctx context.Context,
	loanID uuid.UUID,
	paymentAmount decimal.Decimal,
	makePaymentFn func(loan *entity.Loan, currPaidAmount decimal.Decimal) (payment *entity.LoanPayment, shouldUpdateLoan bool, err error),
) (loan *entity.Loan, newPaidAmount decimal.Decimal, err error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, decimal.Decimal{}, err
	}
	defer func() { err = finishTransaction(err, tx) }()

	loan, err = getLoan(ctx, tx, loanID)
	if err != nil {
		return nil, decimal.Decimal{}, err
	}

	currPaidAmount, err := getLoanPaidAmount(ctx, tx, loanID)
	if err != nil {
		return nil, decimal.Decimal{}, err
	}

	loanPayment, shouldUpdateLoan, err := makePaymentFn(loan, currPaidAmount)
	if err != nil {
		return nil, decimal.Decimal{}, err
	}

	query, args := loanPaymentStruct.InsertInto(loanPaymentsTable, toPostgresLoanPayment(loanPayment)).BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, decimal.Decimal{}, err
	}

	newPaidAmount = currPaidAmount.Add(loanPayment.Amount)

	if shouldUpdateLoan {
		if err = updateLoan(ctx, tx, loan); err != nil {
			return nil, decimal.Decimal{}, err
		}
	}

	return loan, newPaidAmount, nil
}

func getLoan(ctx context.Context, executor executor, loanID uuid.UUID) (*entity.Loan, error) {
	sb := loanStruct.SelectFrom(loansTable)
	query, args := sb.Where(sb.Equal("id", loanID)).BuildWithFlavor(sqlbuilder.PostgreSQL)

	var pgLoan postgresLoan
	err := executor.QueryRowContext(ctx, query, args...).Scan(loanStruct.Addr(&pgLoan)...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return pgLoan.toEntityLoan(), nil
}

func updateLoan(ctx context.Context, executor executor, loan *entity.Loan) error {
	ub := loanStruct.Update(loansTable, toPostgresLoan(loan))
	query, args := ub.Where(ub.Equal("id", loan.ID)).BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := executor.ExecContext(ctx, query, args...)
	return err
}
