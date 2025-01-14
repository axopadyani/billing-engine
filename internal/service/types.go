package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/axopadyani/billing-engine/internal/entity"
)

// LoanStatus represents the status of a loan.
type LoanStatus int

const (
	// LoanStatusOngoing indicates that the loan is still active and payments are ongoing.
	LoanStatusOngoing LoanStatus = iota

	// LoanStatusPaid indicates that the loan has been fully paid off.
	LoanStatusPaid
)

// parseLoanStatus converts an entity.LoanStatus to a service.LoanStatus.
//
// Parameters:
//   - entityStatus: The loan status from the entity package.
//
// Returns:
//   - A LoanStatus corresponding to the input entity status.
func parseLoanStatus(entityStatus entity.LoanStatus) LoanStatus {
	var res LoanStatus
	switch entityStatus {
	case entity.LoanStatusOngoing:
		res = LoanStatusOngoing
	case entity.LoanStatusPaid:
		res = LoanStatusPaid
	}

	return res
}

// Loan represents a loan in the service layer.
type Loan struct {
	ID                   uuid.UUID
	UserID               uuid.UUID
	Amount               decimal.Decimal
	PaymentDurationWeeks int32
	PaymentAmount        decimal.Decimal
	Status               LoanStatus
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// parseLoan converts an entity.Loan to a service.Loan.
//
// Parameters:
//   - entityLoan: A pointer to the loan entity to be converted.
//
// Returns:
//   - A Loan struct populated with data from the entity loan.
//     If entityLoan is nil, an empty Loan struct is returned.
func parseLoan(entityLoan *entity.Loan) Loan {
	if entityLoan == nil {
		return Loan{}
	}

	return Loan{
		ID:                   entityLoan.ID,
		UserID:               entityLoan.UserID,
		Amount:               entityLoan.Amount,
		PaymentDurationWeeks: entityLoan.PaymentDurationWeeks,
		PaymentAmount:        entityLoan.PaymentAmount,
		Status:               parseLoanStatus(entityLoan.Status),
		CreatedAt:            entityLoan.CreatedAt,
		UpdatedAt:            entityLoan.UpdatedAt,
	}
}

// LoanDetail represents detailed information about a loan.
type LoanDetail struct {
	Loan              Loan
	OutstandingAmount decimal.Decimal
	CurrentBillAmount decimal.Decimal
	IsDelinquent      bool
}

// parseLoanDetail creates a LoanDetail struct from individual components.
//
// Parameters:
//   - loan: The base Loan struct.
//   - outstandingAmount: The remaining amount to be paid on the loan.
//   - currentBillAmount: The amount due in the current billing cycle.
//   - isDelinquent: A boolean indicating whether the loan is past due.
//
// Returns:
//   - A LoanDetail struct populated with the provided information.
func parseLoanDetail(loan Loan, outstandingAmount, currentBillAmount decimal.Decimal, isDelinquent bool) LoanDetail {
	return LoanDetail{
		Loan:              loan,
		OutstandingAmount: outstandingAmount,
		CurrentBillAmount: currentBillAmount,
		IsDelinquent:      isDelinquent,
	}
}
