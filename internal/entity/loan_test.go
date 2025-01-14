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

func TestLoanStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status LoanStatus
		want   bool
	}{
		{
			name:   "ongoing",
			status: LoanStatusOngoing,
			want:   true,
		},
		{
			name:   "paid",
			status: LoanStatusPaid,
			want:   true,
		},
		{
			name:   "unknown status",
			status: LoanStatus(-1),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Fatalf("expecting %v, got %v", tt.want, got)
			}
		})
	}
}

func TestLoan_validate(t *testing.T) {
	tests := []struct {
		name      string
		loan      *Loan
		wantError error
	}{
		{
			name: "empty ID",
			loan: &Loan{
				UserID:               uuid.New(),
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 50,
				PaymentAmount:        decimal.NewFromInt(5_500_000),
			},
			wantError: ErrLoanEmptyID,
		},
		{
			name: "empty user ID",
			loan: &Loan{
				ID:                   uuid.New(),
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 50,
				PaymentAmount:        decimal.NewFromInt(5_500_000),
			},
			wantError: ErrLoanEmptyUserID,
		},
		{
			name: "invalid amount",
			loan: &Loan{
				ID:                   uuid.New(),
				UserID:               uuid.New(),
				Amount:               decimal.Zero,
				PaymentDurationWeeks: 50,
				PaymentAmount:        decimal.NewFromInt(5_500_000),
			},
			wantError: ErrLoanInvalidAmount,
		},
		{
			name: "invalid payment duration",
			loan: &Loan{
				ID:                   uuid.New(),
				UserID:               uuid.New(),
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 0,
				PaymentAmount:        decimal.NewFromInt(5_500_000),
			},
			wantError: ErrLoanInvalidPaymentDurationWeeks,
		},
		{
			name: "invalid payment amount",
			loan: &Loan{
				ID:                   uuid.New(),
				UserID:               uuid.New(),
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 50,
				PaymentAmount:        decimal.Zero,
			},
			wantError: ErrLoanInvalidPaymentAmount,
		},
		{
			name: "invalid loan status",
			loan: &Loan{
				ID:                   uuid.New(),
				UserID:               uuid.New(),
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 50,
				PaymentAmount:        decimal.NewFromInt(5_500_000),
				Status:               LoanStatus(-1),
			},
			wantError: ErrLoanInvalidStatus,
		},
		{
			name: "empty created at",
			loan: &Loan{
				ID:                   uuid.New(),
				UserID:               uuid.New(),
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 50,
				PaymentAmount:        decimal.NewFromInt(5_500_000),
				Status:               LoanStatusOngoing,
				CreatedAt:            time.Time{},
			},
			wantError: ErrLoanEmptyCreatedAt,
		},
		{
			name: "empty updated at",
			loan: &Loan{
				ID:                   uuid.New(),
				UserID:               uuid.New(),
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 50,
				PaymentAmount:        decimal.NewFromInt(5_500_000),
				Status:               LoanStatusOngoing,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Time{},
			},
			wantError: ErrLoanEmptyUpdatedAt,
		},
		{
			name: "normal case",
			loan: &Loan{
				ID:                   uuid.New(),
				UserID:               uuid.New(),
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 50,
				PaymentAmount:        decimal.NewFromInt(5_500_000),
				Status:               LoanStatusOngoing,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
			wantError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.loan.validate()
			if !errors.Is(err, test.wantError) {
				t.Fatalf("expecting error to be %v, got %v", test.wantError, err)
			}
		})
	}
}

func TestCreateLoan(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name                 string
		userID               uuid.UUID
		amount               decimal.Decimal
		paymentDurationWeeks int32
		wantLoan             *Loan
		wantErr              error
	}{
		{
			name:                 "empty user ID",
			userID:               uuid.Nil,
			amount:               decimal.NewFromInt(5_000_000),
			paymentDurationWeeks: 50,
			wantLoan:             nil,
			wantErr:              ErrLoanEmptyUserID,
		},
		{
			name:                 "invalid amount",
			userID:               userID,
			amount:               decimal.NewFromInt(0),
			paymentDurationWeeks: 50,
			wantLoan:             nil,
			wantErr:              ErrLoanInvalidAmount,
		},
		{
			name:                 "invalid payment duration",
			userID:               userID,
			amount:               decimal.NewFromInt(5_000_000),
			paymentDurationWeeks: 0,
			wantLoan:             nil,
			wantErr:              ErrLoanInvalidPaymentDurationWeeks,
		},
		{
			name:                 "normal case",
			userID:               userID,
			amount:               decimal.NewFromInt(5_000_000),
			paymentDurationWeeks: 50,
			wantLoan: &Loan{
				UserID:               userID,
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 50,
				PaymentAmount:        decimal.NewFromInt(5_500_000),
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			loan, err := CreateLoan(test.userID, test.amount, test.paymentDurationWeeks)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("expecting error to be %v, got %v", test.wantErr, err)
			}

			if err == nil {
				if diff := cmp.Diff(
					test.wantLoan, loan,
					cmpopts.IgnoreFields(Loan{}, "ID", "CreatedAt", "UpdatedAt"),
				); diff != "" {
					t.Fatalf("loan compare mismatch (-want/+got)\n%s", diff)
				}

				if loan.ID == uuid.Nil {
					t.Fatalf("expecting loan id not to be empty")
				}
				if loan.CreatedAt.IsZero() {
					t.Fatalf("expecting loan created at not to be zero")
				}
				if loan.UpdatedAt.IsZero() {
					t.Fatalf("expecting loan updated at not to be zero")
				}
			}
		})
	}
}

func TestLoan_ValidateLatestLoan(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name       string
		loan       *Loan
		latestLoan *Loan
		wantErr    error
	}{
		{
			name:       "nil latest loan",
			loan:       &Loan{},
			latestLoan: nil,
			wantErr:    nil,
		},
		{
			name:       "nil loan",
			loan:       nil,
			latestLoan: &Loan{},
			wantErr:    nil,
		},
		{
			name: "ongoing loan",
			loan: &Loan{UserID: userID},
			latestLoan: &Loan{
				UserID: userID,
				Status: LoanStatusOngoing,
			},
			wantErr: ErrLoanStillHasOngoingLoan,
		},
		{
			name: "no ongoing loan",
			loan: &Loan{UserID: userID},
			latestLoan: &Loan{
				UserID: userID,
				Status: LoanStatusPaid,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.loan.ValidateLatestLoan(tt.latestLoan)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expecting error to be %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestLoan_OutstandingAmount(t *testing.T) {
	tests := []struct {
		name       string
		loan       *Loan
		paidAmount decimal.Decimal
		wantAmount decimal.Decimal
	}{
		{
			name:       "nil loan",
			loan:       nil,
			paidAmount: decimal.Zero,
			wantAmount: decimal.Zero,
		},
		{
			name:       "full amount outstanding",
			loan:       &Loan{PaymentAmount: decimal.NewFromInt(1000)},
			paidAmount: decimal.Zero,
			wantAmount: decimal.NewFromInt(1000),
		},
		{
			name:       "partial amount paid",
			loan:       &Loan{PaymentAmount: decimal.NewFromInt(1000)},
			paidAmount: decimal.NewFromInt(400),
			wantAmount: decimal.NewFromInt(600),
		},
		{
			name:       "fully paid",
			loan:       &Loan{PaymentAmount: decimal.NewFromInt(1000)},
			paidAmount: decimal.NewFromInt(1000),
			wantAmount: decimal.Zero,
		},
		{
			name:       "overpaid",
			loan:       &Loan{PaymentAmount: decimal.NewFromInt(1000)},
			paidAmount: decimal.NewFromInt(1200),
			wantAmount: decimal.Zero,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.loan.OutstandingAmount(test.paidAmount)
			if !got.Equal(test.wantAmount) {
				t.Fatalf("want %v, got %v", test.wantAmount, got)
			}
		})
	}
}

func TestLoan_IsDelinquent(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name       string
		loan       *Loan
		paidAmount decimal.Decimal
		want       bool
	}{
		{
			name:       "nil loan",
			loan:       nil,
			paidAmount: decimal.Zero,
			want:       false,
		},
		{
			name: "loan is paid",
			loan: &Loan{
				Status:               LoanStatusPaid,
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 7), // 1 week ago
			},
			paidAmount: decimal.NewFromInt(500),
			want:       false,
		},
		{
			name: "paid amount equals payment amount",
			loan: &Loan{
				Status:               LoanStatusOngoing,
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 7), // 1 week ago
			},
			paidAmount: decimal.NewFromInt(1000),
			want:       false,
		},
		{
			name: "delinquent - more than 2 weeks unpaid",
			loan: &Loan{
				Status:               LoanStatusOngoing,
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 21), // 3 weeks ago
			},
			paidAmount: decimal.NewFromInt(0),
			want:       true,
		},
		{
			name: "not delinquent - less than 2 weeks unpaid",
			loan: &Loan{
				Status:               LoanStatusOngoing,
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 7), // 1 week ago
			},
			paidAmount: decimal.NewFromInt(0),
			want:       false,
		},
		{
			name: "exactly 2 weeks unpaid",
			loan: &Loan{
				Status:               LoanStatusOngoing,
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 14), // 2 weeks ago
			},
			paidAmount: decimal.NewFromInt(0),
			want:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.loan.IsDelinquent(now, test.paidAmount)
			if got != test.want {
				t.Fatalf("expecting delinquency to be %t, got %t", test.want, got)
			}
		})
	}
}

func TestLoan_CurrentBillAmount(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name               string
		loan               *Loan
		paidAmount         decimal.Decimal
		expectedBillAmount decimal.Decimal
	}{
		{
			name:               "nil loan",
			loan:               nil,
			paidAmount:         decimal.Zero,
			expectedBillAmount: decimal.Zero,
		},
		{
			name: "first week, no payment",
			loan: &Loan{
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 7), // now is loan week 1
			},
			paidAmount:         decimal.Zero,
			expectedBillAmount: decimal.NewFromInt(100), // weekly payment amount
		},
		{
			name: "mid-duration, partial payment",
			loan: &Loan{
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 21), // now is loan week 3
			},
			paidAmount:         decimal.NewFromInt(100),
			expectedBillAmount: decimal.NewFromInt(200),
		},
		{
			name: "after duration, no payment",
			loan: &Loan{
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 84), // now is loan week 12
			},
			paidAmount:         decimal.Zero,
			expectedBillAmount: decimal.NewFromInt(1000),
		},
		{
			name: "after duration, partial payment",
			loan: &Loan{
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 84), // now is loan week 12
			},
			paidAmount:         decimal.NewFromInt(600),
			expectedBillAmount: decimal.NewFromInt(400),
		},
		{
			name: "fully paid for the week",
			loan: &Loan{
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 21), // now is loan week 3
			},
			paidAmount:         decimal.NewFromInt(300),
			expectedBillAmount: decimal.Zero,
		},
		{
			name: "fully paid",
			loan: &Loan{
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 21), // now is loan week 3
			},
			paidAmount:         decimal.NewFromInt(1000),
			expectedBillAmount: decimal.Zero,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.loan.CurrentBillAmount(now, test.paidAmount)
			if !got.Equal(test.expectedBillAmount) {
				t.Fatalf("expected bill amount to be %s, got %s", test.expectedBillAmount, got)
			}
		})
	}
}

func TestLoan_MakePayment(t *testing.T) {
	now := time.Now().UTC()
	loanID := uuid.New()

	tests := []struct {
		name            string
		loan            *Loan
		paidAmount      decimal.Decimal
		paymentAmount   decimal.Decimal
		wantLoanPayment *LoanPayment
		wantUpdateLoan  bool
		wantErr         error
	}{
		{
			name:            "nil loan",
			loan:            nil,
			paidAmount:      decimal.Zero,
			paymentAmount:   decimal.NewFromInt(100),
			wantLoanPayment: nil,
			wantUpdateLoan:  false,
			wantErr:         ErrLoanNotFound,
		},
		{
			name: "current week already paid",
			loan: &Loan{
				ID:                   loanID,
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 7), // now is loan week 1
			},
			paidAmount:      decimal.NewFromInt(100), // week 1 already paid
			paymentAmount:   decimal.NewFromInt(100),
			wantLoanPayment: nil,
			wantUpdateLoan:  false,
			wantErr:         ErrLoanCurrentWeekAlreadyPaid,
		},
		{
			name: "payment amount does not match bill amount",
			loan: &Loan{
				ID:                   loanID,
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 7), // now is loan week 1
			},
			paidAmount:      decimal.Zero,
			paymentAmount:   decimal.NewFromInt(50), // should be 100 every week
			wantLoanPayment: nil,
			wantUpdateLoan:  false,
			wantErr:         ErrLoanNotExactPaymentAmount,
		},
		{
			name: "not the last week's payment, should not update loan",
			loan: &Loan{
				ID:                   loanID,
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 7), // now is loan week 1
			},
			paidAmount:    decimal.Zero,
			paymentAmount: decimal.NewFromInt(100),
			wantLoanPayment: &LoanPayment{
				LoanID: loanID,
				Amount: decimal.NewFromInt(100),
			},
			wantUpdateLoan: false,
			wantErr:        nil,
		},
		{
			name: "last week's payment, should update loan",
			loan: &Loan{
				ID:                   loanID,
				PaymentAmount:        decimal.NewFromInt(1000),
				PaymentDurationWeeks: 10,
				CreatedAt:            now.Add(-time.Hour * 24 * 70), // now is loan week 10
			},
			paidAmount:    decimal.NewFromInt(800),
			paymentAmount: decimal.NewFromInt(200),
			wantLoanPayment: &LoanPayment{
				LoanID: loanID,
				Amount: decimal.NewFromInt(200),
			},
			wantUpdateLoan: true,
			wantErr:        nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			loanPayment, shouldUpdateLoan, err := test.loan.MakePayment(now, test.paidAmount, test.paymentAmount)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("expecting error %v, got %v", test.wantErr, err)
				return
			}

			if err == nil {
				if diff := cmp.Diff(
					test.wantLoanPayment, loanPayment,
					cmpopts.IgnoreFields(LoanPayment{}, "ID", "CreatedAt", "UpdatedAt"),
				); diff != "" {
					t.Fatalf("LoanPayment missmatch (-want +got):\n%s", diff)
				}

				if loanPayment.ID == uuid.Nil {
					t.Fatalf("expecting loanPayment.ID to be non-zero")
				}

				if loanPayment.CreatedAt.IsZero() {
					t.Errorf("MakePayment() loanPayment.CreatedAt should not be zero")
				}

				if shouldUpdateLoan != test.wantUpdateLoan {
					t.Errorf("MakePayment() shouldUpdateLoan = %v, want %v", shouldUpdateLoan, test.wantUpdateLoan)
				}

				if test.wantUpdateLoan {
					if test.loan.Status != LoanStatusPaid {
						t.Errorf("MakePayment() loan status should be LoanStatusPaid")
					}
					if test.loan.UpdatedAt.IsZero() {
						t.Errorf("MakePayment() loan UpdatedAt should not be zero")
					}
				}
			}
		})
	}
}

func TestLoan_weeklyPaymentAmount(t *testing.T) {
	tests := []struct {
		name string
		loan *Loan
		want decimal.Decimal
	}{
		{
			name: "nil loan",
			loan: nil,
			want: decimal.Zero,
		},
		{
			name: "even division",
			loan: &Loan{
				PaymentAmount:        decimal.NewFromInt(10000),
				PaymentDurationWeeks: 10,
			},
			want: decimal.NewFromInt(1000),
		},
		{
			name: "uneven division",
			loan: &Loan{
				PaymentAmount:        decimal.NewFromInt(10000),
				PaymentDurationWeeks: 3,
			},
			want: decimal.NewFromInt(3333),
		},
		{
			name: "one week duration",
			loan: &Loan{
				PaymentAmount:        decimal.NewFromInt(5000),
				PaymentDurationWeeks: 1,
			},
			want: decimal.NewFromInt(5000),
		},
		{
			name: "zero payment amount",
			loan: &Loan{
				PaymentAmount:        decimal.Zero,
				PaymentDurationWeeks: 5,
			},
			want: decimal.Zero,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.loan.weeklyPaymentAmount()
			if !got.Equal(test.want) {
				t.Fatalf("expecting weekly payment amount to be %s, got %s", test.want, got)
			}
		})
	}
}

func TestLoan_currentWeek(t *testing.T) {
	tests := []struct {
		name string
		loan *Loan
		now  time.Time
		want int32
	}{
		{
			name: "nil loan",
			loan: nil,
			now:  time.Now(),
			want: 0,
		},
		{
			name: "normal case",
			loan: &Loan{
				CreatedAt: time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC), // Monday
			},
			now:  time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), // 2 weeks later
			want: 2,
		},
		{
			name: "same week",
			loan: &Loan{
				CreatedAt: time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC), // Monday
			},
			now:  time.Date(2023, 5, 7, 23, 59, 59, 0, time.UTC), // Sunday of the same week
			want: 0,
		},
		{
			name: "created on Sunday",
			loan: &Loan{
				CreatedAt: time.Date(2023, 5, 7, 23, 59, 59, 0, time.UTC), // Sunday
			},
			now:  time.Date(2023, 5, 8, 0, 0, 0, 0, time.UTC), // Monday of the next week
			want: 1,
		},
		{
			name: "now in the past",
			loan: &Loan{
				CreatedAt: time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC), // Monday
			},
			now:  time.Date(2023, 4, 30, 23, 59, 59, 0, time.UTC), // Monday of the previous week
			want: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.loan.currentWeek(test.now)
			if got != test.want {
				t.Fatalf("expecting current week to be %d, got %d", test.want, got)
			}
		})
	}
}
