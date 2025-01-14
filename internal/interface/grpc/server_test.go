package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/axopadyani/billing-engine/internal/service"
	mock "github.com/axopadyani/billing-engine/internal/test/mock/service"
	v1 "github.com/axopadyani/billing-engine/proto/v1"
)

var testTimeout = 10 * time.Second

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockService(ctrl)
	server := NewServer(mockSvc)
	if server == nil {
		t.Error("expecting server to be created")
	}
}

func TestServer_CreateLoan(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	mockRes := service.Loan{
		ID:                   uuid.New(),
		UserID:               uuid.New(),
		Amount:               decimal.NewFromInt(5_000_000),
		PaymentDurationWeeks: 5,
		PaymentAmount:        decimal.NewFromInt(5_500_000),
		Status:               service.LoanStatusOngoing,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	tests := []struct {
		name      string
		setupMock func(*mock.MockService)
		request   *v1.CreateLoanRequest
		wantErr   *status.Status
	}{
		{
			name:      "invalid user id",
			setupMock: nil,
			request: &v1.CreateLoanRequest{
				UserId:               "invalid",
				Amount:               mockRes.Amount.String(),
				PaymentDurationWeeks: mockRes.PaymentDurationWeeks,
			},
			wantErr: status.New(codes.InvalidArgument, "invalid user id"),
		},
		{
			name:      "invalid amount",
			setupMock: nil,
			request: &v1.CreateLoanRequest{
				UserId:               mockRes.UserID.String(),
				Amount:               "invalid",
				PaymentDurationWeeks: mockRes.PaymentDurationWeeks,
			},
			wantErr: status.New(codes.InvalidArgument, "invalid amount"),
		},
		{
			name: "loan service error",
			setupMock: func(mockSvc *mock.MockService) {
				mockSvc.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(service.Loan{}, service.UnexpectedError)
			},
			request: &v1.CreateLoanRequest{
				UserId:               mockRes.UserID.String(),
				Amount:               mockRes.Amount.String(),
				PaymentDurationWeeks: mockRes.PaymentDurationWeeks,
			},
			wantErr: status.New(codes.Internal, service.UnexpectedError.Error()),
		},
		{
			name: "normal case",
			setupMock: func(mockSvc *mock.MockService) {
				mockSvc.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(mockRes, nil)
			},
			request: &v1.CreateLoanRequest{
				UserId:               mockRes.UserID.String(),
				Amount:               mockRes.Amount.String(),
				PaymentDurationWeeks: mockRes.PaymentDurationWeeks,
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mock.NewMockService(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockSvc)
			}

			server := NewServer(mockSvc)

			_, err := server.CreateLoan(ctx, test.request)
			if err != nil {
				statusErr, ok := status.FromError(err)
				if !ok {
					t.Fatalf("unexpected error: %v", err)
				}

				if test.wantErr.Message() != statusErr.Message() {
					t.Fatalf("expecting error message %q, got %q", test.wantErr.Message(), statusErr.Message())
				}
				if test.wantErr.Code() != statusErr.Code() {
					t.Fatalf("expecting error code %v, got %v", test.wantErr.Code(), statusErr.Code())
				}
			} else if err == nil && test.wantErr != nil {
				t.Fatal("expecting error not to be nil")
			}
		})
	}
}

func TestServer_GetCurrentLoan(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	mockLoan := service.Loan{
		ID:                   uuid.New(),
		UserID:               uuid.New(),
		Amount:               decimal.NewFromInt(5_000_000),
		PaymentDurationWeeks: 5,
		PaymentAmount:        decimal.NewFromInt(5_500_000),
		Status:               service.LoanStatusOngoing,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	tests := []struct {
		name      string
		setupMock func(mockSvc *mock.MockService)
		req       *v1.GetCurrentLoanRequest
		wantErr   *status.Status
	}{
		{
			name:      "invalid user id",
			setupMock: nil,
			req:       &v1.GetCurrentLoanRequest{UserId: "invalid"},
			wantErr:   status.New(codes.InvalidArgument, "invalid user id"),
		},
		{
			name: "service error",
			setupMock: func(mockSvc *mock.MockService) {
				mockSvc.EXPECT().GetCurrentLoan(gomock.Any(), gomock.Any()).Return(service.LoanDetail{}, service.UnexpectedError)
			},
			req:     &v1.GetCurrentLoanRequest{UserId: uuid.NewString()},
			wantErr: status.New(codes.Internal, service.UnexpectedError.Error()),
		},
		{
			name: "normal case",
			setupMock: func(mockSvc *mock.MockService) {
				mockSvc.EXPECT().GetCurrentLoan(gomock.Any(), gomock.Any()).Return(
					service.LoanDetail{
						Loan:              mockLoan,
						OutstandingAmount: mockLoan.PaymentAmount,
						CurrentBillAmount: decimal.Zero,
						IsDelinquent:      false,
					},
					nil,
				)
			},
			req:     &v1.GetCurrentLoanRequest{UserId: uuid.NewString()},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mock.NewMockService(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockSvc)
			}

			server := NewServer(mockSvc)
			_, err := server.GetCurrentLoan(ctx, test.req)
			if err != nil {
				statusErr, ok := status.FromError(err)
				if !ok {
					t.Fatalf("unexpected error: %v", err)
				}

				if test.wantErr.Message() != statusErr.Message() {
					t.Fatalf("expecting error message %q, got %q", test.wantErr.Message(), statusErr.Message())
				}
				if test.wantErr.Code() != statusErr.Code() {
					t.Fatalf("expecting error code %v, got %v", test.wantErr.Code(), statusErr.Code())
				}
			} else if err == nil && test.wantErr != nil {
				t.Fatal("expecting error not to be nil")
			}
		})
	}
}

func TestServer_MakePayment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	mockLoanDetail := service.LoanDetail{
		Loan: service.Loan{
			ID:                   uuid.New(),
			UserID:               uuid.New(),
			Amount:               decimal.NewFromInt(5_000_000),
			PaymentDurationWeeks: 5,
			PaymentAmount:        decimal.NewFromInt(5_500_000),
			Status:               service.LoanStatusOngoing,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		},
		OutstandingAmount: decimal.NewFromInt(5_500_000),
		CurrentBillAmount: decimal.NewFromInt(1_100_000),
		IsDelinquent:      false,
	}

	tests := []struct {
		name      string
		setupMock func(*mock.MockService)
		req       *v1.MakePaymentRequest
		wantErr   *status.Status
	}{
		{
			name:      "invalid loan id",
			setupMock: nil,
			req:       &v1.MakePaymentRequest{LoanId: "invalid", PaymentAmount: "1000000"},
			wantErr:   status.New(codes.InvalidArgument, "invalid user id"),
		},
		{
			name:      "invalid payment amount",
			setupMock: nil,
			req:       &v1.MakePaymentRequest{LoanId: uuid.New().String(), PaymentAmount: "invalid"},
			wantErr:   status.New(codes.InvalidArgument, "invalid payment amount"),
		},
		{
			name: "service error",
			setupMock: func(mockSvc *mock.MockService) {
				mockSvc.EXPECT().MakePayment(gomock.Any(), gomock.Any()).Return(service.LoanDetail{}, service.UnexpectedError)
			},
			req:     &v1.MakePaymentRequest{LoanId: uuid.New().String(), PaymentAmount: "1000000"},
			wantErr: status.New(codes.Internal, service.UnexpectedError.Error()),
		},
		{
			name: "normal case",
			setupMock: func(mockSvc *mock.MockService) {
				mockSvc.EXPECT().MakePayment(gomock.Any(), gomock.Any()).Return(mockLoanDetail, nil)
			},
			req:     &v1.MakePaymentRequest{LoanId: uuid.New().String(), PaymentAmount: "1000000"},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mock.NewMockService(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockSvc)
			}

			server := NewServer(mockSvc)
			_, err := server.MakePayment(ctx, test.req)
			if err != nil {
				statusErr, ok := status.FromError(err)
				if !ok {
					t.Fatalf("unexpected error: %v", err)
				}

				if test.wantErr.Message() != statusErr.Message() {
					t.Fatalf("expecting error message %q, got %q", test.wantErr.Message(), statusErr.Message())
				}
				if test.wantErr.Code() != statusErr.Code() {
					t.Fatalf("expecting error code %v, got %v", test.wantErr.Code(), statusErr.Code())
				}
			} else if err == nil && test.wantErr != nil {
				t.Fatal("expecting error not to be nil")
			}
		})
	}
}
