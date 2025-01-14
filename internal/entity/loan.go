package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/axopadyani/billing-engine/internal/common/businesserror"
)

const delinquencyThresholdWeeks = 2 // Number of unpaid weeks to be considered delinquent

var (
	ErrLoanEmptyID                     = businesserror.New("loan id cannot be empty", businesserror.KindBadRequest)
	ErrLoanEmptyUserID                 = businesserror.New("loan user id cannot be empty", businesserror.KindBadRequest)
	ErrLoanInvalidAmount               = businesserror.New("loan amount must be greater than zero", businesserror.KindBadRequest)
	ErrLoanInvalidPaymentDurationWeeks = businesserror.New("loan payment duration must be at least 1 week", businesserror.KindBadRequest)
	ErrLoanInvalidPaymentAmount        = businesserror.New("loan payment amount must be greater than zero", businesserror.KindBadRequest)
	ErrLoanInvalidStatus               = businesserror.New("invalid loan status", businesserror.KindBadRequest)
	ErrLoanEmptyCreatedAt              = businesserror.New("created at cannot be empty", businesserror.KindBadRequest)
	ErrLoanEmptyUpdatedAt              = businesserror.New("updated at cannot be empty", businesserror.KindBadRequest)
	ErrLoanStillHasOngoingLoan         = businesserror.New("user still has ongoing loan", businesserror.KindUnprocessableEntity)
	ErrLoanNotFound                    = businesserror.New("loan not found", businesserror.KindNotFound)
	ErrLoanCurrentWeekAlreadyPaid      = businesserror.New("current week is already paid", businesserror.KindUnprocessableEntity)
	ErrLoanNotExactPaymentAmount       = businesserror.New("loan payment amount does not match billing amount", businesserror.KindUnprocessableEntity)

	interestRate = decimal.NewFromFloat(0.1)
)

// LoanStatus represents the current state of a loan.
type LoanStatus int

const (
	// LoanStatusOngoing indicates that the loan is still active and payments are ongoing.
	LoanStatusOngoing LoanStatus = iota

	// LoanStatusPaid indicates that the loan has been fully paid off.
	LoanStatusPaid
)

// IsValid checks if the LoanStatus is a valid status.
//
// This method determines whether the LoanStatus is one of the predefined valid statuses.
//
// Returns:
//   - bool: true if the status is either LoanStatusOngoing or LoanStatusPaid, false otherwise.
func (s LoanStatus) IsValid() bool {
	return s == LoanStatusOngoing || s == LoanStatusPaid
}

// Loan represents a loan entity in the system.
//
// It contains all the necessary information about a loan, including its
// unique identifier, the user it belongs to, the loan amount, payment duration,
// total payment amount (including interest), current status, and timestamps.
type Loan struct {
	// ID is the unique identifier for the loan.
	ID uuid.UUID

	// UserID is the unique identifier of the user who took the loan.
	UserID uuid.UUID

	// Amount is the principal amount of the loan.
	Amount decimal.Decimal

	// PaymentDurationWeeks is the duration of the loan in weeks.
	PaymentDurationWeeks int32

	// PaymentAmount is the total amount to be paid, including interest.
	PaymentAmount decimal.Decimal

	// Status represents the current state of the loan (e.g., ongoing, paid).
	Status LoanStatus

	// CreatedAt is the timestamp when the loan was created.
	CreatedAt time.Time

	// UpdatedAt is the timestamp when the loan was last updated.
	UpdatedAt time.Time
}

// validate checks if the Loan instance is valid by verifying all its fields.
// It ensures that:
// - The loan ID is not empty
// - The user ID is not empty
// - The loan amount is greater than zero
// - The payment duration is at least 1 week
// - The payment amount is greater than zero
// - The loan status is valid
// - The creation and update timestamps are not zero
//
// Returns:
//   - error: An error if any validation check fails, nil if the loan is valid.
func (l *Loan) validate() error {
	if l.ID == uuid.Nil {
		return ErrLoanEmptyID
	}

	if l.UserID == uuid.Nil {
		return ErrLoanEmptyUserID
	}

	if l.Amount.LessThanOrEqual(decimal.Zero) {
		return ErrLoanInvalidAmount
	}

	if l.PaymentDurationWeeks <= 0 {
		return ErrLoanInvalidPaymentDurationWeeks
	}

	if l.PaymentAmount.LessThanOrEqual(decimal.Zero) {
		return ErrLoanInvalidPaymentAmount
	}

	if !l.Status.IsValid() {
		return ErrLoanInvalidStatus
	}

	if l.CreatedAt.IsZero() {
		return ErrLoanEmptyCreatedAt
	}

	if l.UpdatedAt.IsZero() {
		return ErrLoanEmptyUpdatedAt
	}

	return nil
}

// CreateLoan creates a new loan for a user with the specified details.
//
// Parameters:
//   - userID: The unique identifier of the user taking the loan.
//   - amount: The principal amount of the loan.
//   - paymentDurationWeeks: The duration of the loan in weeks.
//
// Returns:
//   - *Loan: A pointer to the newly created Loan instance if successful.
//   - error: An error if the loan creation fails, nil otherwise.
//
// The function generates a new UUID for the loan, calculates the total payment amount
// (including interest), and sets the initial status to ongoing. It also performs
// validation on the created loan instance before returning.
func CreateLoan(userID uuid.UUID, amount decimal.Decimal, paymentDurationWeeks int32) (*Loan, error) {
	loanID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	loan := &Loan{
		ID:                   loanID,
		UserID:               userID,
		Amount:               amount,
		PaymentDurationWeeks: paymentDurationWeeks,
		PaymentAmount:        amount.Add(amount.Mul(interestRate).RoundUp(0)),
		Status:               LoanStatusOngoing,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err = loan.validate(); err != nil {
		return nil, err
	}

	return loan, nil
}

// ValidateLatestLoan checks if the user associated with this loan has any ongoing loans.
// It compares the current loan with the latest loan to determine if a new loan can be created.
//
// Parameters:
//   - latestLoan: A pointer to the most recent Loan instance for the user.
//
// Returns:
//   - An error if the user has an ongoing loan, nil otherwise.
//
// The function returns ErrLoanStillHasOngoingLoan if the user associated with this loan
// already has an ongoing loan.
func (l *Loan) ValidateLatestLoan(latestLoan *Loan) error {
	if l != nil && latestLoan != nil && l.UserID == latestLoan.UserID && latestLoan.Status == LoanStatusOngoing {
		return ErrLoanStillHasOngoingLoan
	}

	return nil
}

// OutstandingAmount calculates the remaining amount to be paid on the loan.
//
// This method subtracts the paid amount from the total payment amount of the loan.
// If the result is negative, it returns zero, ensuring the outstanding amount is never negative.
//
// Parameters:
//   - paidAmount: The decimal.Decimal amount that has already been paid towards the loan.
//
// Returns:
//   - decimal.Decimal: The outstanding amount to be paid.
func (l *Loan) OutstandingAmount(paidAmount decimal.Decimal) decimal.Decimal {
	if l == nil {
		return decimal.Zero
	}

	outstandingAmount := l.PaymentAmount.Sub(paidAmount)
	if outstandingAmount.IsNegative() {
		outstandingAmount = decimal.Zero
	}

	return outstandingAmount
}

// IsDelinquent determines if the loan is considered delinquent based on the current date and paid amount.
//
// A loan is considered delinquent if the number of unpaid weeks exceeds the delinquencyThresholdWeeks.
//
// Parameters:
//   - now: The current time used to calculate the billing amount.
//   - paidAmount: The total amount that has been paid towards the loan so far.
//
// Returns:
//   - bool: true if the loan is delinquent, false otherwise.
func (l *Loan) IsDelinquent(now time.Time, paidAmount decimal.Decimal) bool {
	if l == nil {
		return false
	}

	if l.Status == LoanStatusPaid || l.PaymentAmount.Equal(paidAmount) {
		return false
	}

	billAmount := l.CurrentBillAmount(now, paidAmount)
	unpaidWeeks := billAmount.Div(l.weeklyPaymentAmount()).Round(0).IntPart()
	return unpaidWeeks > delinquencyThresholdWeeks
}

// CurrentBillAmount calculates the current bill amount for the loan based on the current date and paid amount.
//
// This method determines the amount that should be billed to the user at the current point in time,
// taking into account the loan's payment schedule and any amounts already paid.
//
// Parameters:
//   - now: The current time used to calculate the billing amount.
//   - paidAmount: The total amount that has been paid towards the loan so far.
//
// Returns:
//   - decimal.Decimal: The current bill amount. This will be zero if the loan is fully paid.
func (l *Loan) CurrentBillAmount(now time.Time, paidAmount decimal.Decimal) decimal.Decimal {
	if l == nil {
		return decimal.Zero
	}

	paymentObligation := l.PaymentAmount

	// cap the amount to the total payment amount
	if currentWeek := l.currentWeek(now); currentWeek < l.PaymentDurationWeeks {
		paymentObligation = l.weeklyPaymentAmount().Mul(decimal.NewFromInt32(currentWeek))
	}

	billAmount := paymentObligation.Sub(paidAmount)
	if billAmount.IsNegative() {
		billAmount = decimal.Zero
	}

	return billAmount
}

// MakePayment processes a payment for the loan and updates its status if necessary.
//
// This method checks if the payment amount matches the current bill amount, creates a new
// loan payment instance, and determines if the loan status should be updated to paid.
//
// Parameters:
//   - now: The current time used to calculate the current bill amount.
//   - paidAmount: The total amount already paid towards the loan before this payment.
//   - paymentAmount: The amount being paid in this transaction.
//
// Returns:
//   - loanPayment: The newly created LoanPayment instance.
//   - shouldUpdateLoan: A boolean indicating whether any changes being made to the loan instance.
//   - err: An error if the payment process fails, nil otherwise. Possible errors include:
//     ErrLoanNotFound, ErrLoanCurrentWeekAlreadyPaid, ErrLoanNotExactPaymentAmount.
func (l *Loan) MakePayment(now time.Time, paidAmount, paymentAmount decimal.Decimal) (loanPayment *LoanPayment, shouldUpdateLoan bool, err error) {
	if l == nil {
		return nil, false, ErrLoanNotFound
	}

	billAmount := l.CurrentBillAmount(now, paidAmount)
	if billAmount.IsZero() {
		return nil, false, ErrLoanCurrentWeekAlreadyPaid
	}
	if !billAmount.Equal(paymentAmount) {
		return nil, false, ErrLoanNotExactPaymentAmount
	}

	loanPayment, err = CreateLoanPayment(l.ID, paymentAmount)
	if err != nil {
		return nil, false, err
	}

	shouldUpdateLoan = false
	if paidAmount.Add(paymentAmount).Equal(l.PaymentAmount) {
		l.Status = LoanStatusPaid
		l.UpdatedAt = time.Now().UTC()
		shouldUpdateLoan = true
	}

	return loanPayment, shouldUpdateLoan, nil
}

// weeklyPaymentAmount calculates the weekly payment amount for the loan.
//
// This method computes the amount to be paid each week by dividing the total payment amount
// by the number of weeks in the loan duration. The result is rounded down to the nearest
// whole number.
//
// Returns:
//   - decimal.Decimal: The amount that should be paid in weekly-basis.
func (l *Loan) weeklyPaymentAmount() decimal.Decimal {
	if l == nil {
		return decimal.Zero
	}

	return l.PaymentAmount.Div(decimal.NewFromInt32(l.PaymentDurationWeeks)).RoundDown(0)
}

// currentWeek calculates the number of weeks that have passed since the loan was created.
//
// This method determines the current week of the loan by calculating the difference
// between the given time and the beginning of the week when the loan was created.
// It accounts for partial weeks and ensures the count starts from the Monday of the
// creation week.
//
// Parameters:
//   - now: The current time to calculate the week difference from.
//
// Returns:
//   - int32: The number of weeks that have passed since the loan was created.
//     Returns 0 if the loan is nil.
func (l *Loan) currentWeek(now time.Time) int32 {
	if l == nil {
		return 0
	}

	createdAt := l.CreatedAt.UTC()

	// get the Monday's date of the loan's creation week
	weekday := int(createdAt.Weekday() - 1)
	if weekday < 0 {
		weekday += 7
	}
	beginningOfWeek := createdAt.AddDate(0, 0, -weekday)

	beginningOfWeek = time.Date(beginningOfWeek.Year(), beginningOfWeek.Month(), beginningOfWeek.Day(), 0, 0, 0, 0, time.UTC)
	currentWeek := int32(now.Sub(beginningOfWeek).Hours() / (24 * 7))
	return currentWeek
}
