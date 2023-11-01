// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ilya-rusyanov/gophermart/internal/usecases (interfaces: AccrualService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entities "github.com/ilya-rusyanov/gophermart/internal/entities"
)

// MockAccrualService is a mock of AccrualService interface.
type MockAccrualService struct {
	ctrl     *gomock.Controller
	recorder *MockAccrualServiceMockRecorder
}

// MockAccrualServiceMockRecorder is the mock recorder for MockAccrualService.
type MockAccrualServiceMockRecorder struct {
	mock *MockAccrualService
}

// NewMockAccrualService creates a new mock instance.
func NewMockAccrualService(ctrl *gomock.Controller) *MockAccrualService {
	mock := &MockAccrualService{ctrl: ctrl}
	mock.recorder = &MockAccrualServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccrualService) EXPECT() *MockAccrualServiceMockRecorder {
	return m.recorder
}

// GetStateOfOrder mocks base method.
func (m *MockAccrualService) GetStateOfOrder(arg0 context.Context, arg1 entities.OrderID) (entities.OrderStatus, entities.Currency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateOfOrder", arg0, arg1)
	ret0, _ := ret[0].(entities.OrderStatus)
	ret1, _ := ret[1].(entities.Currency)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetStateOfOrder indicates an expected call of GetStateOfOrder.
func (mr *MockAccrualServiceMockRecorder) GetStateOfOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateOfOrder", reflect.TypeOf((*MockAccrualService)(nil).GetStateOfOrder), arg0, arg1)
}
