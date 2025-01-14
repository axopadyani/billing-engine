package grpc

import (
	"context"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/axopadyani/billing-engine/internal/service"
	v1 "github.com/axopadyani/billing-engine/proto/v1"
)

// Server represents the gRPC server for the Billing Engine.
type Server struct {
	v1.UnimplementedBillingEngineServer
	svc service.Service
}

// NewServer creates a new instance of the Billing Engine gRPC server.
//
// Parameters:
//   - svc: The service implementation for handling business logic.
//
// Returns:
//   - The newly created Server instance.
func NewServer(svc service.Service) *Server {
	return &Server{
		svc: svc,
	}
}

// CreateLoan handles the creation of a new loan for a user.
//
// Parameters:
//   - ctx: The context for the request.
//   - in: The v1.CreateLoanRequest protobuf message.
//
// Returns:
//   - The created loan as v1.Loan protobuf message.
//   - An error if the loan creation fails or input is invalid.
func (s *Server) CreateLoan(ctx context.Context, in *v1.CreateLoanRequest) (*v1.Loan, error) {
	userID, err := uuid.Parse(in.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	amount, err := decimal.NewFromString(in.GetAmount())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid amount")
	}

	res, err := s.svc.CreateLoan(ctx, service.CreateLoanCommand{
		UserID:               userID,
		Amount:               amount,
		PaymentDurationWeeks: in.GetPaymentDurationWeeks(),
	})
	if err != nil {
		return nil, toGrpcError(err)
	}

	return parseLoan(res), nil
}

// GetCurrentLoan retrieves the current loan details for a user.
//
// Parameters:
//   - ctx: The context for the request.
//   - in: The v1.GetCurrentLoanRequest protobuf message.
//
// Returns:
//   - The loan detail as v1.LoanDetail protobuf message.
//   - An error if retrieval fails or input is invalid.
func (s *Server) GetCurrentLoan(ctx context.Context, in *v1.GetCurrentLoanRequest) (*v1.LoanDetail, error) {
	userID, err := uuid.Parse(in.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	res, err := s.svc.GetCurrentLoan(ctx, service.GetCurrentLoanQuery{UserID: userID})
	if err != nil {
		return nil, toGrpcError(err)
	}

	return parseLoanDetail(res), nil
}

// MakePayment processes a payment for a specific loan.
//
// Parameters:
//   - ctx: The context for the request.
//   - in: The v1.MakePaymentRequest protobuf message.
//
// Returns:
//   - The updated loan details as v1.LoanDetail protobuf message.
//   - An error if the payment fails or input is invalid.
func (s *Server) MakePayment(ctx context.Context, in *v1.MakePaymentRequest) (*v1.LoanDetail, error) {
	loanID, err := uuid.Parse(in.GetLoanId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	paymentAmount, err := decimal.NewFromString(in.GetPaymentAmount())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payment amount")
	}

	res, err := s.svc.MakePayment(ctx, service.MakePaymentCommand{
		LoanID:        loanID,
		PaymentAmount: paymentAmount,
	})
	if err != nil {
		return nil, toGrpcError(err)
	}

	return parseLoanDetail(res), nil
}

// Serve starts the gRPC server and begins listening for incoming requests.
//
// Parameters:
//   - listener: The net.Listener to use for accepting connections.
//
// This function will block to serve requests until it is stopped or encounters a fatal error.
func (s *Server) Serve(listener net.Listener) {
	grpcServer := grpc.NewServer()
	v1.RegisterBillingEngineServer(grpcServer, s)
	reflection.Register(grpcServer)

	log.Printf("server listening on %s", listener.Addr().String())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
