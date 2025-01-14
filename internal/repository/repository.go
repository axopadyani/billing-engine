package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/axopadyani/billing-engine/internal/entity"
)

// Repository defines the interface for repository operations related to loans and payments.
//
//go:generate mockgen -package repository -source=repository.go -destination=../test/mock/repository/mock_repository.go
type Repository interface {
    // CreateLoan creates a new loan in the repository.
    //
    // Parameters:
    //   - ctx: The context for the operation.
    //   - loan: A pointer to the Loan entity to be created.
    //   - validateFn: A function to validate the loan before creation.
    //
    // Returns:
    //   An error if the creation fails, nil otherwise.
    CreateLoan(ctx context.Context, loan *entity.Loan, validateFn func(latestLoan *entity.Loan) error) error

    // GetLatestLoan retrieves the most recent loan for a given user.
    //
    // Parameters:
    //   - ctx: The context for the operation.
    //   - userID: The UUID of the user whose latest loan is to be retrieved.
    //
    // Returns:
    //   A pointer to the latest Loan entity and an error if the retrieval fails.
    GetLatestLoan(ctx context.Context, userID uuid.UUID) (*entity.Loan, error)

    // GetLoanPaidAmount retrieves the total amount paid for a specific loan.
    //
    // Parameters:
    //   - ctx: The context for the operation.
    //   - loanID: The UUID of the loan for which to get the paid amount.
    //
    // Returns:
    //   The paid amount as a decimal.Decimal and an error if the retrieval fails.
    GetLoanPaidAmount(ctx context.Context, loanID uuid.UUID) (decimal.Decimal, error)

    // MakePayment processes a payment for a loan.
    //
    // Parameters:
    //   - ctx: The context for the operation.
    //   - loanID: The UUID of the loan for which the payment is being made.
    //   - paymentAmount: The amount of the payment as a decimal.Decimal.
    //   - makePaymentFn: A function to process the payment and determine if the loan should be updated.
    //
    // Returns:
    //   The updated Loan entity, the new total paid amount, and an error if the payment processing fails.
    MakePayment(
        ctx context.Context,
        loanID uuid.UUID,
        paymentAmount decimal.Decimal,
        makePaymentFn func(loan *entity.Loan, currPaidAmount decimal.Decimal) (payment *entity.LoanPayment, shouldUpdateLoan bool, err error),
    ) (loan *entity.Loan, newPaidAmount decimal.Decimal, err error)
}
