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

func TestImpl_CreateLoan(t *testing.T) {
	userID := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	tests := []struct {
		name      string
		setupMock func(mockRepo *repository.MockRepository)
		cmd       CreateLoanCommand
		wantErr   error
	}{
		{
			name:      "validation error",
			setupMock: nil,
			cmd: CreateLoanCommand{
				UserID:               uuid.Nil,
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 5,
			},
			wantErr: entity.ErrLoanEmptyUserID,
		},
		{
			name: "repo business error",
			setupMock: func(mockRepo *repository.MockRepository) {
				mockRepo.EXPECT().
					CreateLoan(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(entity.ErrLoanStillHasOngoingLoan)
			},
			cmd: CreateLoanCommand{
				UserID:               userID,
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 5,
			},
			wantErr: entity.ErrLoanStillHasOngoingLoan,
		},
		{
			name: "repo unexpected error",
			setupMock: func(mockRepo *repository.MockRepository) {
				mockRepo.EXPECT().
					CreateLoan(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("unknown error"))
			},
			cmd: CreateLoanCommand{
				UserID:               userID,
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 5,
			},
			wantErr: UnexpectedError,
		},
		{
			name: "normal case",
			setupMock: func(mockRepo *repository.MockRepository) {
				mockRepo.EXPECT().
					CreateLoan(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			cmd: CreateLoanCommand{
				UserID:               userID,
				Amount:               decimal.NewFromInt(5_000_000),
				PaymentDurationWeeks: 5,
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockRepository(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockRepo)
			}

			s := NewService(mockRepo)

			_, err := s.CreateLoan(ctx, test.cmd)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("expecting error to be %v, got %v", test.wantErr, err)
			}
		})
	}
}
