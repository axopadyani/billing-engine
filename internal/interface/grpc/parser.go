package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/axopadyani/billing-engine/internal/service"
	v1 "github.com/axopadyani/billing-engine/proto/v1"
)

// parseLoan converts a service.Loan to a v1.Loan protobuf message.
//
// Parameters:
//   - loan: A service.Loan struct containing the loan information.
//
// Returns:
//   - *v1.Loan: A pointer to a v1.Loan struct with the converted loan data.
func parseLoan(loan service.Loan) *v1.Loan {
	return &v1.Loan{
		Id:                   loan.ID.String(),
		UserId:               loan.UserID.String(),
		Amount:               loan.Amount.String(),
		PaymentDurationWeeks: loan.PaymentDurationWeeks,
		PaymentAmount:        loan.PaymentAmount.String(),
		Status:               parseLoanStatus(loan.Status),
		CreatedAt:            timestamppb.New(loan.CreatedAt),
		UpdatedAt:            timestamppb.New(loan.UpdatedAt),
	}
}

// parseLoanStatus converts a service.LoanStatus to a v1.LoanStatus protobuf enum.
//
// Parameters:
//   - status: A service.LoanStatus representing the internal loan status.
//
// Returns:
//   - v1.LoanStatus: The corresponding v1.LoanStatus enum value.
func parseLoanStatus(status service.LoanStatus) v1.LoanStatus {
	var res v1.LoanStatus
	switch status {
	case service.LoanStatusOngoing:
		res = v1.LoanStatus_ONGOING
	case service.LoanStatusPaid:
		res = v1.LoanStatus_PAID
	}

	return res
}

// parseLoanDetail converts a service.LoanDetail to a v1.LoanDetail protobuf message.
//
// Parameters:
//   - loanDetail: A service.LoanDetail struct containing the loan detail information.
//
// Returns:
//   - *v1.LoanDetail: A pointer to a v1.LoanDetail struct with the converted loan detail data.
func parseLoanDetail(loanDetail service.LoanDetail) *v1.LoanDetail {
	return &v1.LoanDetail{
		Loan:              parseLoan(loanDetail.Loan),
		OutstandingAmount: loanDetail.OutstandingAmount.String(),
		CurrentBillAmount: loanDetail.CurrentBillAmount.String(),
		IsDelinquent:      loanDetail.IsDelinquent,
	}
}
