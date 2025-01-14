package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/axopadyani/billing-engine/internal/entity"
)

// CreateLoanCommand represents the input data required to create a new loan.
type CreateLoanCommand struct {
	// UserID is the unique identifier of the user requesting the loan.
	UserID uuid.UUID

	// Amount is the decimal representation of the loan amount.
	Amount decimal.Decimal

	// PaymentDurationWeeks is the duration of the loan repayment period in weeks.
	PaymentDurationWeeks int32
}

// CreateLoan creates a new loan for a user based on the provided command.
//
// It first creates a loan entity, then validates it against the latest loan (if any),
// and finally persists it in the repository.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for cancellation or passing values.
//   - in: A CreateLoanCommand struct containing the necessary information to create a loan.
//
// Returns:
//   - Loan: A Loan struct representing the created loan if successful.
//   - error: An error if the loan creation fails, or nil if successful.
func (s *Impl) CreateLoan(ctx context.Context, in CreateLoanCommand) (Loan, error) {
	loan, err := entity.CreateLoan(in.UserID, in.Amount, in.PaymentDurationWeeks)
	if err != nil {
		return Loan{}, ensureBusinessError(err)
	}

	err = s.repo.CreateLoan(ctx, loan, func(latestLoan *entity.Loan) error {
		return loan.ValidateLatestLoan(latestLoan)
	})

	if err != nil {
		return Loan{}, ensureBusinessError(err)
	}

	return parseLoan(loan), nil
}
