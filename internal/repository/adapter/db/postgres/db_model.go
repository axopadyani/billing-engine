package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/shopspring/decimal"

	"github.com/axopadyani/billing-engine/internal/entity"
)

const (
	loansTable        = "loans"
	loanPaymentsTable = "loan_payments"
)

// postgresLoan represents a loan record in the PostgreSQL database.
type postgresLoan struct {
	ID                   uuid.UUID       `db:"id"`
	UserID               uuid.UUID       `db:"user_id"`
	Amount               decimal.Decimal `db:"amount"`
	PaymentDurationWeeks int32           `db:"payment_duration_weeks"`
	PaymentAmount        decimal.Decimal `db:"payment_amount"`
	Status               int             `db:"status"`
	CreatedAt            time.Time       `db:"created_at"`
	UpdatedAt            time.Time       `db:"updated_at"`
}

var loanStruct = sqlbuilder.NewStruct(new(postgresLoan))

func toPostgresLoan(loan *entity.Loan) *postgresLoan {
	return &postgresLoan{
		ID:                   loan.ID,
		UserID:               loan.UserID,
		Amount:               loan.Amount,
		PaymentDurationWeeks: loan.PaymentDurationWeeks,
		PaymentAmount:        loan.PaymentAmount,
		Status:               int(loan.Status),
		CreatedAt:            loan.CreatedAt,
		UpdatedAt:            loan.UpdatedAt,
	}
}

func (l postgresLoan) toEntityLoan() *entity.Loan {
	return &entity.Loan{
		ID:                   l.ID,
		UserID:               l.UserID,
		Amount:               l.Amount,
		PaymentDurationWeeks: l.PaymentDurationWeeks,
		PaymentAmount:        l.PaymentAmount,
		Status:               entity.LoanStatus(l.Status),
		CreatedAt:            l.CreatedAt,
		UpdatedAt:            l.UpdatedAt,
	}
}

// postgresLoanPayment represents a loan payment record in the PostgreSQL database.
type postgresLoanPayment struct {
	ID        uuid.UUID       `db:"id"`
	LoanID    uuid.UUID       `db:"loan_id"`
	Amount    decimal.Decimal `db:"amount"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

var loanPaymentStruct = sqlbuilder.NewStruct(new(postgresLoanPayment))

func toPostgresLoanPayment(loanPayment *entity.LoanPayment) *postgresLoanPayment {
	return &postgresLoanPayment{
		ID:        loanPayment.ID,
		LoanID:    loanPayment.LoanID,
		Amount:    loanPayment.Amount,
		CreatedAt: loanPayment.CreatedAt,
		UpdatedAt: loanPayment.UpdatedAt,
	}
}
