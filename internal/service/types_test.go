package service

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/axopadyani/billing-engine/internal/entity"
)

func TestParseLoanStatus(t *testing.T) {
	tests := []struct {
		name         string
		entityStatus entity.LoanStatus
		want         LoanStatus
	}{
		{
			name:         "ongoing",
			entityStatus: entity.LoanStatusOngoing,
			want:         LoanStatusOngoing,
		},
		{
			name:         "pain",
			entityStatus: entity.LoanStatusPaid,
			want:         LoanStatusPaid,
		},
		{
			name:         "unknown",
			entityStatus: entity.LoanStatus(999),
			want:         LoanStatus(0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := parseLoanStatus(test.entityStatus); got != test.want {
				t.Fatalf("expecting %v, got %v", test.want, got)
			}
		})
	}
}

func TestParseLoan(t *testing.T) {
	mockLoan, err := entity.CreateLoan(uuid.New(), decimal.NewFromInt(5_000_000), 5)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		entityLoan *entity.Loan
		want       Loan
	}{
		{
			name:       "nil entity loan",
			entityLoan: nil,
			want:       Loan{},
		},
		{
			name:       "normal case",
			entityLoan: mockLoan,
			want: Loan{
				ID:                   mockLoan.ID,
				UserID:               mockLoan.UserID,
				Amount:               mockLoan.Amount,
				PaymentDurationWeeks: mockLoan.PaymentDurationWeeks,
				PaymentAmount:        mockLoan.PaymentAmount,
				Status:               parseLoanStatus(mockLoan.Status),
				CreatedAt:            mockLoan.CreatedAt,
				UpdatedAt:            mockLoan.UpdatedAt,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseLoan(test.entityLoan)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Fatalf("parseLoan() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseLoanDetail(t *testing.T) {
	mockLoan := Loan{
		ID:                   uuid.New(),
		UserID:               uuid.New(),
		Amount:               decimal.NewFromInt(1000000),
		PaymentDurationWeeks: 10,
		PaymentAmount:        decimal.NewFromInt(110000),
		Status:               LoanStatusOngoing,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	tests := []struct {
		name              string
		loan              Loan
		outstandingAmount decimal.Decimal
		currentBillAmount decimal.Decimal
		isDelinquent      bool
		wantLoanDetail    LoanDetail
	}{
		{
			name:              "normal case",
			loan:              mockLoan,
			outstandingAmount: decimal.NewFromInt(500000),
			currentBillAmount: decimal.NewFromInt(110000),
			isDelinquent:      false,
			wantLoanDetail: LoanDetail{
				Loan:              mockLoan,
				OutstandingAmount: decimal.NewFromInt(500000),
				CurrentBillAmount: decimal.NewFromInt(110000),
				IsDelinquent:      false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseLoanDetail(test.loan, test.outstandingAmount, test.currentBillAmount, test.isDelinquent)
			if diff := cmp.Diff(test.wantLoanDetail, got); diff != "" {
				t.Fatalf("parseLoanDetail() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
