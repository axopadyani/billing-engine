package service

import (
	"context"
	"errors"

	"github.com/axopadyani/billing-engine/internal/common/businesserror"
	"github.com/axopadyani/billing-engine/internal/repository"
)

var UnexpectedError = businesserror.New("unexpected error, please try again", businesserror.KindInternal)

// Service defines the interface for the billing engine operations.
//
//go:generate mockgen -package service -source=service.go -destination=../test/mock/service/mock_service.go
type Service interface {
	// CreateLoan creates a new loan.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - cmd: The CreateLoanCommand containing the loan creation details.
	//
	// Returns:
	//   - Loan: The created loan information.
	//   - error: An error if the operation fails, or nil if successful.
	CreateLoan(ctx context.Context, cmd CreateLoanCommand) (Loan, error)

	// GetCurrentLoan retrieves the current loan details.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - query: The GetCurrentLoanQuery containing the query parameters.
	//
	// Returns:
	//   - LoanDetail: The details of the current loan.
	//   - error: An error if the operation fails, or nil if successful.
	GetCurrentLoan(ctx context.Context, query GetCurrentLoanQuery) (LoanDetail, error)

	// MakePayment processes a payment for a loan.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - cmd: The MakePaymentCommand containing the payment details.
	//
	// Returns:
	//   - LoanDetail: The updated loan details after the payment.
	//   - error: An error if the operation fails, or nil if successful.
	MakePayment(ctx context.Context, cmd MakePaymentCommand) (LoanDetail, error)
}

// Impl represents the implementation of the Service interface.
type Impl struct {
	// repo is the repository interface used for data storage and retrieval operations.
	repo repository.Repository
}

// NewService creates and returns a new instance of the Service implementation.
//
// It initializes the Impl struct with the provided repository.
//
// Parameters:
//   - repo: A repository.Repository interface implementation used for data storage and retrieval operations.
//
// Returns:
//   - *Impl: The newly created Impl struct, which implements the Service interface.
func NewService(repo repository.Repository) *Impl {
	return &Impl{repo: repo}
}

// ensureBusinessError wraps non-business errors with a generic UnexpectedError.
// It ensures that all errors returned from the service are of type BusinessError.
//
// Parameters:
//   - err: The error to be checked.
//
// Returns:
//   - error: Either the original error if it's already a BusinessError,
//     nil if the input is nil, UnexpectedError for other error types.
func ensureBusinessError(err error) error {
	if err == nil {
		return nil
	}

	var businessErr *businesserror.BusinessError
	if errors.As(err, &businessErr) {
		return businessErr
	}
	return UnexpectedError
}
