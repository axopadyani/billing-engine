package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/axopadyani/billing-engine/internal/entity"
	"github.com/axopadyani/billing-engine/internal/test/mock/repository"
)

func TestImpl_MakePayment(t *testing.T) {
	ctx := context.Background()

	mockLoan, err := entity.CreateLoan(uuid.New(), decimal.NewFromInt(5_000_000), 5)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name      string
		setupMock func(mockRepo *repository.MockRepository)
		cmd       MakePaymentCommand
		wantErr   error
	}{
		{
			name: "normal case",
			setupMock: func(mockRepo *repository.MockRepository) {
				mockRepo.EXPECT().MakePayment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mockLoan, decimal.NewFromInt(2000), nil)
			},
			cmd: MakePaymentCommand{
				LoanID:        mockLoan.ID,
				PaymentAmount: decimal.NewFromInt(1000),
			},
			wantErr: nil,
		},
		{
			name: "repository expected error",
			setupMock: func(mockRepo *repository.MockRepository) {
				mockRepo.EXPECT().MakePayment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, decimal.Zero, entity.ErrLoanNotFound)
			},
			wantErr: entity.ErrLoanNotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockRepository(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockRepo)
			}

			svc := NewService(mockRepo)

			_, err := svc.MakePayment(ctx, test.cmd)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("expecting error to be %v, got %v", test.wantErr, err)
			}
		})
	}
}
