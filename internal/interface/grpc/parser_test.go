package grpc

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/axopadyani/billing-engine/internal/service"
	v1 "github.com/axopadyani/billing-engine/proto/v1"
)

func TestParseLoan(t *testing.T) {
	now := time.Now()
	input := service.Loan{
		ID:                   uuid.New(),
		UserID:               uuid.New(),
		Amount:               decimal.NewFromInt(5000000),
		PaymentDurationWeeks: 50,
		PaymentAmount:        decimal.NewFromInt(5500000),
		Status:               service.LoanStatusOngoing,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	want := &v1.Loan{
		Id:                   input.ID.String(),
		UserId:               input.UserID.String(),
		Amount:               "5000000",
		PaymentDurationWeeks: 50,
		PaymentAmount:        "5500000",
		Status:               v1.LoanStatus_ONGOING,
		CreatedAt:            timestamppb.New(now),
		UpdatedAt:            timestamppb.New(now),
	}

	got := parseLoan(input)

	if diff := cmp.Diff(
		want, got,
		cmpopts.IgnoreUnexported(v1.Loan{}, timestamppb.Timestamp{}),
	); diff != "" {
		t.Fatalf("parseLoan() mismatch (-want +got):\n%s", diff)
	}
}

func TestParseLoanStatus(t *testing.T) {
	tests := []struct {
		name   string
		status service.LoanStatus
		want   v1.LoanStatus
	}{
		{
			name:   "ongoing",
			status: service.LoanStatusOngoing,
			want:   v1.LoanStatus_ONGOING,
		},
		{
			name:   "paid",
			status: service.LoanStatusPaid,
			want:   v1.LoanStatus_PAID,
		},
		{
			name:   "unknown status",
			status: service.LoanStatus(-1),
			want:   v1.LoanStatus(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLoanStatus(tt.status)
			if got != tt.want {
				t.Fatalf("parseLoanStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLoanDetail(t *testing.T) {
	now := time.Now()
	input := service.LoanDetail{
		Loan: service.Loan{
			ID:                   uuid.New(),
			UserID:               uuid.New(),
			Amount:               decimal.NewFromInt(5000000),
			PaymentDurationWeeks: 50,
			PaymentAmount:        decimal.NewFromInt(5500000),
			Status:               service.LoanStatusOngoing,
			CreatedAt:            now,
			UpdatedAt:            now,
		},
		OutstandingAmount: decimal.NewFromInt(3000000),
		CurrentBillAmount: decimal.NewFromInt(100000),
		IsDelinquent:      false,
	}

	want := &v1.LoanDetail{
		Loan: &v1.Loan{
			Id:                   input.Loan.ID.String(),
			UserId:               input.Loan.UserID.String(),
			Amount:               "5000000",
			PaymentDurationWeeks: 50,
			PaymentAmount:        "5500000",
			Status:               v1.LoanStatus_ONGOING,
			CreatedAt:            timestamppb.New(now),
			UpdatedAt:            timestamppb.New(now),
		},
		OutstandingAmount: "3000000",
		CurrentBillAmount: "100000",
		IsDelinquent:      false,
	}

	got := parseLoanDetail(input)

	if diff := cmp.Diff(
		want, got,
		cmpopts.IgnoreUnexported(v1.LoanDetail{}, v1.Loan{}, timestamppb.Timestamp{}),
	); diff != "" {
		t.Fatalf("parseLoanDetail() mismatch (-want +got):\n%s", diff)
	}
}
