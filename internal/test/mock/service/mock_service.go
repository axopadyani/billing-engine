// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	service "github.com/axopadyani/billing-engine/internal/service"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateLoan mocks base method.
func (m *MockService) CreateLoan(ctx context.Context, cmd service.CreateLoanCommand) (service.Loan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLoan", ctx, cmd)
	ret0, _ := ret[0].(service.Loan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLoan indicates an expected call of CreateLoan.
func (mr *MockServiceMockRecorder) CreateLoan(ctx, cmd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLoan", reflect.TypeOf((*MockService)(nil).CreateLoan), ctx, cmd)
}

// GetCurrentLoan mocks base method.
func (m *MockService) GetCurrentLoan(ctx context.Context, query service.GetCurrentLoanQuery) (service.LoanDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentLoan", ctx, query)
	ret0, _ := ret[0].(service.LoanDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentLoan indicates an expected call of GetCurrentLoan.
func (mr *MockServiceMockRecorder) GetCurrentLoan(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentLoan", reflect.TypeOf((*MockService)(nil).GetCurrentLoan), ctx, query)
}

// MakePayment mocks base method.
func (m *MockService) MakePayment(ctx context.Context, cmd service.MakePaymentCommand) (service.LoanDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakePayment", ctx, cmd)
	ret0, _ := ret[0].(service.LoanDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakePayment indicates an expected call of MakePayment.
func (mr *MockServiceMockRecorder) MakePayment(ctx, cmd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakePayment", reflect.TypeOf((*MockService)(nil).MakePayment), ctx, cmd)
}