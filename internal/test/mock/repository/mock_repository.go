// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	entity "github.com/axopadyani/billing-engine/internal/entity"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	decimal "github.com/shopspring/decimal"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateLoan mocks base method.
func (m *MockRepository) CreateLoan(ctx context.Context, loan *entity.Loan, validateFn func(*entity.Loan) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLoan", ctx, loan, validateFn)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLoan indicates an expected call of CreateLoan.
func (mr *MockRepositoryMockRecorder) CreateLoan(ctx, loan, validateFn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLoan", reflect.TypeOf((*MockRepository)(nil).CreateLoan), ctx, loan, validateFn)
}

// GetLatestLoan mocks base method.
func (m *MockRepository) GetLatestLoan(ctx context.Context, userID uuid.UUID) (*entity.Loan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestLoan", ctx, userID)
	ret0, _ := ret[0].(*entity.Loan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestLoan indicates an expected call of GetLatestLoan.
func (mr *MockRepositoryMockRecorder) GetLatestLoan(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestLoan", reflect.TypeOf((*MockRepository)(nil).GetLatestLoan), ctx, userID)
}

// GetLoanPaidAmount mocks base method.
func (m *MockRepository) GetLoanPaidAmount(ctx context.Context, loanID uuid.UUID) (decimal.Decimal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoanPaidAmount", ctx, loanID)
	ret0, _ := ret[0].(decimal.Decimal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoanPaidAmount indicates an expected call of GetLoanPaidAmount.
func (mr *MockRepositoryMockRecorder) GetLoanPaidAmount(ctx, loanID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoanPaidAmount", reflect.TypeOf((*MockRepository)(nil).GetLoanPaidAmount), ctx, loanID)
}

// MakePayment mocks base method.
func (m *MockRepository) MakePayment(ctx context.Context, loanID uuid.UUID, paymentAmount decimal.Decimal, makePaymentFn func(*entity.Loan, decimal.Decimal) (*entity.LoanPayment, bool, error)) (*entity.Loan, decimal.Decimal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakePayment", ctx, loanID, paymentAmount, makePaymentFn)
	ret0, _ := ret[0].(*entity.Loan)
	ret1, _ := ret[1].(decimal.Decimal)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// MakePayment indicates an expected call of MakePayment.
func (mr *MockRepositoryMockRecorder) MakePayment(ctx, loanID, paymentAmount, makePaymentFn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakePayment", reflect.TypeOf((*MockRepository)(nil).MakePayment), ctx, loanID, paymentAmount, makePaymentFn)
}
