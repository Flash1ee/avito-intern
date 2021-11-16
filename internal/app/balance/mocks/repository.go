// Code generated by MockGen. DO NOT EDIT.
// Source: avito-intern/internal/app/balance (interfaces: Repository)

// Package mock_balance is a generated GoMock package.
package mock_balance

import (
	models "avito-intern/internal/app/balance/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// BalanceRepository is a mock of Repository interface.
type BalanceRepository struct {
	ctrl     *gomock.Controller
	recorder *BalanceRepositoryMockRecorder
}

// BalanceRepositoryMockRecorder is the mock recorder for BalanceRepository.
type BalanceRepositoryMockRecorder struct {
	mock *BalanceRepository
}

// NewBalanceRepository creates a new mock instance.
func NewBalanceRepository(ctrl *gomock.Controller) *BalanceRepository {
	mock := &BalanceRepository{ctrl: ctrl}
	mock.recorder = &BalanceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *BalanceRepository) EXPECT() *BalanceRepositoryMockRecorder {
	return m.recorder
}

// AddBalance mocks base method.
func (m *BalanceRepository) AddBalance(arg0 int64, arg1 float64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBalance", arg0, arg1)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddBalance indicates an expected call of AddBalance.
func (mr *BalanceRepositoryMockRecorder) AddBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBalance", reflect.TypeOf((*BalanceRepository)(nil).AddBalance), arg0, arg1)
}

// CreateAccount mocks base method.
func (m *BalanceRepository) CreateAccount(arg0 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *BalanceRepositoryMockRecorder) CreateAccount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*BalanceRepository)(nil).CreateAccount), arg0)
}

// CreateTransfer mocks base method.
func (m *BalanceRepository) CreateTransfer(arg0, arg1 int64, arg2 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransfer", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTransfer indicates an expected call of CreateTransfer.
func (mr *BalanceRepositoryMockRecorder) CreateTransfer(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransfer", reflect.TypeOf((*BalanceRepository)(nil).CreateTransfer), arg0, arg1, arg2)
}

// FindUserByID mocks base method.
func (m *BalanceRepository) FindUserByID(arg0 int64) (*models.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByID", arg0)
	ret0, _ := ret[0].(*models.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByID indicates an expected call of FindUserByID.
func (mr *BalanceRepositoryMockRecorder) FindUserByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByID", reflect.TypeOf((*BalanceRepository)(nil).FindUserByID), arg0)
}