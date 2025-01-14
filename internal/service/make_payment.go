package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/axopadyani/billing-engine/internal/entity"
)

// MakePaymentCommand represents the input data required to process a loan payment.
type MakePaymentCommand struct {
	// LoanID is the unique identifier of the loan for which the payment is being made.
	LoanID uuid.UUID

	// PaymentAmount is the decimal amount of the payment being made towards the loan.
	PaymentAmount decimal.Decimal
}

// MakePayment processes a payment for a loan.
//
// It updates the loan's payment status, calculates the new paid amount,
// and returns the updated loan details.
//
// Parameters:
//   - ctx: The context for the operation.
//   - in: A MakePaymentCommand struct containing the necessary information to process the payment.
//
// Returns:
//   - LoanDetail: A struct containing the updated loan information.
//   - error: An error if the payment process fails, or nil if successful.
func (s *Impl) MakePayment(ctx context.Context, in MakePaymentCommand) (LoanDetail, error) {
	now := time.Now().UTC()

	loan, newPaidAmount, err := s.repo.MakePayment(
		ctx, in.LoanID, in.PaymentAmount,
		func(loan *entity.Loan, currPaidAmount decimal.Decimal) (payment *entity.LoanPayment, shouldUpdateLoan bool, err error) {
			return loan.MakePayment(now, currPaidAmount, in.PaymentAmount)
		},
	)

	if err != nil {
		return LoanDetail{}, ensureBusinessError(err)
	}

	return parseLoanDetail(
		parseLoan(loan),
		loan.OutstandingAmount(newPaidAmount),
		loan.CurrentBillAmount(now, newPaidAmount),
		loan.IsDelinquent(now, newPaidAmount),
	), nil
}
