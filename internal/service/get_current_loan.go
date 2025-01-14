package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/axopadyani/billing-engine/internal/entity"
)

// GetCurrentLoanQuery represents a query to retrieve the current loan for a user.
type GetCurrentLoanQuery struct {
	// UserID is the unique identifier of the user whose current loan is being queried.
	UserID uuid.UUID
}

// GetCurrentLoan retrieves the current loan details for a given user.
//
// It fetches the latest loan for the user, calculates the outstanding amount,
// current bill amount, and checks if the loan is delinquent.
//
// Parameters:
//   - ctx: The context for the function call, which can be used for cancellation or passing request-scoped values.
//   - in: A GetCurrentLoanQuery struct containing the necessary information to retrieve the current loan details.
//
// Returns:
//   - LoanDetail: A struct containing the detailed information about the current loan.
//   - error: An error if any occurred during the process. It returns entity.ErrLoanNotFound if no ongoing loan is found.
func (s *Impl) GetCurrentLoan(ctx context.Context, in GetCurrentLoanQuery) (LoanDetail, error) {
	loan, err := s.repo.GetLatestLoan(ctx, in.UserID)
	if err != nil {
		return LoanDetail{}, ensureBusinessError(err)
	}
	if loan == nil || loan.Status != entity.LoanStatusOngoing {
		return LoanDetail{}, entity.ErrLoanNotFound
	}

	now := time.Now()
	paidAmount, err := s.repo.GetLoanPaidAmount(ctx, loan.ID)
	if err != nil {
		return LoanDetail{}, ensureBusinessError(err)
	}

	return parseLoanDetail(
		parseLoan(loan),
		loan.OutstandingAmount(paidAmount),
		loan.CurrentBillAmount(now, paidAmount),
		loan.IsDelinquent(now, paidAmount),
	), nil
}
