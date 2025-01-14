package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestLoanPayment_validate(t *testing.T) {
	validID := uuid.New()
	validAmount := decimal.NewFromInt(1000)
	validTime := time.Now().UTC()

	tests := []struct {
		name    string
		payment *LoanPayment
		wantErr error
	}{
		{
			name: "valid payment",
			payment: &LoanPayment{
				ID:        validID,
				LoanID:    validID,
				Amount:    validAmount,
				CreatedAt: validTime,
				UpdatedAt: validTime,
			},
			wantErr: nil,
		},
		{
			name: "empty ID",
			payment: &LoanPayment{
				ID:        uuid.Nil,
				LoanID:    validID,
				Amount:    validAmount,
				CreatedAt: validTime,
				UpdatedAt: validTime,
			},
			wantErr: ErrLoanPaymentEmptyID,
		},
		{
			name: "empty loan id",
			payment: &LoanPayment{
				ID:        validID,
				LoanID:    uuid.Nil,
				Amount:    validAmount,
				CreatedAt: validTime,
				UpdatedAt: validTime,
			},
			wantErr: ErrLoanPaymentEmptyLoanID,
		},
		{
			name: "invalid amount",
			payment: &LoanPayment{
				ID:        validID,
				LoanID:    validID,
				Amount:    decimal.Zero,
				CreatedAt: validTime,
				UpdatedAt: validTime,
			},
			wantErr: ErrLoanPaymentInvalidAmount,
		},
		{
			name: "empty created at",
			payment: &LoanPayment{
				ID:        validID,
				LoanID:    validID,
				Amount:    validAmount,
				CreatedAt: time.Time{},
				UpdatedAt: validTime,
			},
			wantErr: ErrLoanPaymentEmptyCreatedAt,
		},
		{
			name: "empty updated at",
			payment: &LoanPayment{
				ID:        validID,
				LoanID:    validID,
				Amount:    validAmount,
				CreatedAt: validTime,
				UpdatedAt: time.Time{},
			},
			wantErr: ErrLoanPaymentEmptyUpdatedAt,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.payment.validate()
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("expecting error to be %v, got %v", test.wantErr, err)
			}
		})
	}
}

func TestCreateLoanPayment(t *testing.T) {
	validLoanID := uuid.New()
	validAmount := decimal.NewFromInt(1000)

	tests := []struct {
		name    string
		loanID  uuid.UUID
		amount  decimal.Decimal
		wantRes *LoanPayment
		wantErr error
	}{
		{
			name:   "valid payment",
			loanID: validLoanID,
			amount: validAmount,
			wantRes: &LoanPayment{
				LoanID: validLoanID,
				Amount: validAmount,
			},
			wantErr: nil,
		},
		{
			name:    "validation error",
			loanID:  uuid.Nil,
			amount:  validAmount,
			wantRes: nil,
			wantErr: ErrLoanPaymentEmptyLoanID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := CreateLoanPayment(test.loanID, test.amount)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("expecting error to be %v, got %v", test.wantErr, err)
			}

			if err == nil {
				if diff := cmp.Diff(
					test.wantRes, res,
					cmpopts.IgnoreFields(LoanPayment{}, "ID", "CreatedAt", "UpdatedAt"),
				); diff != "" {
					t.Fatalf("LoanPayment missmatch (-want +got):\n%s", diff)
				}

				if res.ID == uuid.Nil {
					t.Fatal("expecting loan payment id to be non-zero")
				}

				if res.CreatedAt.IsZero() {
					t.Fatal("expecting loan payment created at to be non-zero")
				}

				if res.UpdatedAt.IsZero() {
					t.Fatal("expecting loan payment updated at to be non-zero")
				}
			}
		})
	}
}
