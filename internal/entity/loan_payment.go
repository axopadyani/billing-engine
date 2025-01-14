package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/axopadyani/billing-engine/internal/common/businesserror"
)

var (
	ErrLoanPaymentEmptyID        = businesserror.New("loan payment id cannot be empty", businesserror.KindBadRequest)
	ErrLoanPaymentEmptyLoanID    = businesserror.New("loan payment loan id cannot be empty", businesserror.KindBadRequest)
	ErrLoanPaymentInvalidAmount  = businesserror.New("loan payment amount must be greater than zero", businesserror.KindBadRequest)
	ErrLoanPaymentEmptyCreatedAt = businesserror.New("created at cannot be empty", businesserror.KindBadRequest)
	ErrLoanPaymentEmptyUpdatedAt = businesserror.New("updated at cannot be empty", businesserror.KindBadRequest)
)

// LoanPayment represents a payment made towards a loan.
type LoanPayment struct {
    // ID is the unique identifier for the loan payment.
    ID uuid.UUID

    // LoanID is the unique identifier of the loan associated with this payment.
    LoanID uuid.UUID

    // Amount is the monetary value of the payment.
    Amount decimal.Decimal

    // CreatedAt is the timestamp when the payment record was created.
    CreatedAt time.Time

    // UpdatedAt is the timestamp when the payment record was last updated.
    UpdatedAt time.Time
}

// CreateLoanPayment creates a new LoanPayment instance with the given loan ID and amount.
// It generates a new UUID for the payment, sets the creation and update times to the current UTC time,
// and validates the payment before returning it.
//
// Parameters:
//   - loanID: A UUID representing the ID of the loan associated with this payment.
//   - amount: A decimal.Decimal value representing the amount of the payment.
//
// Returns:
//   - *LoanPayment: The newly created and validated LoanPayment instance.
//   - error: An error if there was a problem creating the UUID or if the payment fails validation.
func CreateLoanPayment(loanID uuid.UUID, amount decimal.Decimal) (*LoanPayment, error) {
    paymentID, err := uuid.NewV7()
    if err != nil {
        return nil, err
    }

    now := time.Now().UTC()
    payment := &LoanPayment{
        ID:        paymentID,
        LoanID:    loanID,
        Amount:    amount,
        CreatedAt: now,
        UpdatedAt: now,
    }

    if err = payment.validate(); err != nil {
        return nil, err
    }

    return payment, nil
}

// validate checks the LoanPayment struct for validity.
//
// It performs the following checks:
//   - Ensures the ID is not empty (nil UUID)
//   - Ensures the LoanID is not empty (nil UUID)
//   - Verifies that the Amount is greater than zero
//   - Checks that CreatedAt is not a zero time
//   - Checks that UpdatedAt is not a zero time
//
// Returns:
//   - error: nil if the LoanPayment is valid, otherwise returns a specific error
//     indicating which validation check failed.
func (lp *LoanPayment) validate() error {
    if lp.ID == uuid.Nil {
        return ErrLoanPaymentEmptyID
    }

    if lp.LoanID == uuid.Nil {
        return ErrLoanPaymentEmptyLoanID
    }

    if lp.Amount.LessThanOrEqual(decimal.Zero) {
        return ErrLoanPaymentInvalidAmount
    }

    if lp.CreatedAt.IsZero() {
        return ErrLoanPaymentEmptyCreatedAt
    }

    if lp.UpdatedAt.IsZero() {
        return ErrLoanPaymentEmptyUpdatedAt
    }

    return nil
}
