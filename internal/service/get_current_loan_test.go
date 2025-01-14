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

func TestImpl_GetCurrentLoan(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	ongoingLoan, err := entity.CreateLoan(uuid.New(), decimal.NewFromInt(5_000_000), 5)
	if err != nil {
		t.Fatal(err)
	}

	paidLoan, err := entity.CreateLoan(uuid.New(), decimal.NewFromInt(5_000_000), 5)
	if err != nil {
		t.Fatal(err)
	}
	paidLoan.Status = entity.LoanStatusPaid

	tests := []struct {
		name      string
		cmd       GetCurrentLoanQuery
		setupMock func(*repository.MockRepository)
		wantErr   error
	}{
		{
			name: "loan not found",
			cmd:  GetCurrentLoanQuery{UserID: uuid.New()},
			setupMock: func(mockRepo *repository.MockRepository) {
				mockRepo.EXPECT().GetLatestLoan(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			wantErr: entity.ErrLoanNotFound,
		},
		{
			name: "get loan paid amount unexpected error",
			cmd:  GetCurrentLoanQuery{UserID: uuid.New()},
			setupMock: func(mockRepo *repository.MockRepository) {
				mockRepo.EXPECT().GetLatestLoan(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error"))
			},
			wantErr: UnexpectedError,
		},
		{
			name: "paid loan",
			cmd:  GetCurrentLoanQuery{UserID: uuid.New()},
			setupMock: func(mockRepo *repository.MockRepository) {
				mockRepo.EXPECT().GetLatestLoan(gomock.Any(), gomock.Any()).Return(paidLoan, nil)
			},
			wantErr: entity.ErrLoanNotFound,
		},
		{
			name: "get loan unexpected error",
			cmd:  GetCurrentLoanQuery{UserID: uuid.New()},
			setupMock: func(mockRepo *repository.MockRepository) {
				mockRepo.EXPECT().GetLatestLoan(gomock.Any(), gomock.Any()).Return(ongoingLoan, nil)
				mockRepo.EXPECT().GetLoanPaidAmount(gomock.Any(), gomock.Any()).Return(decimal.Zero, errors.New("unexpected error"))
			},
			wantErr: UnexpectedError,
		},
		{
			name: "normal case",
			cmd:  GetCurrentLoanQuery{UserID: uuid.New()},
			setupMock: func(mockRepo *repository.MockRepository) {
				mockRepo.EXPECT().GetLatestLoan(gomock.Any(), gomock.Any()).Return(ongoingLoan, nil)
				mockRepo.EXPECT().GetLoanPaidAmount(gomock.Any(), gomock.Any()).Return(decimal.Zero, nil)
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockRepository(ctrl)
			test.setupMock(mockRepo)

			s := NewService(mockRepo)

			_, err := s.GetCurrentLoan(ctx, test.cmd)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("expecting error to be %v, got %v", test.wantErr, err)
				return
			}
		})
	}
}
