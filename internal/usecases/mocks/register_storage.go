// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ilya-rusyanov/gophermart/internal/usecases (interfaces: RegisterStorage)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	entities "github.com/ilya-rusyanov/gophermart/internal/entities"
)

// MockRegisterStorage is a mock of RegisterStorage interface.
type MockRegisterStorage struct {
	ctrl     *gomock.Controller
	recorder *MockRegisterStorageMockRecorder
}

// MockRegisterStorageMockRecorder is the mock recorder for MockRegisterStorage.
type MockRegisterStorageMockRecorder struct {
	mock *MockRegisterStorage
}

// NewMockRegisterStorage creates a new mock instance.
func NewMockRegisterStorage(ctrl *gomock.Controller) *MockRegisterStorage {
	mock := &MockRegisterStorage{ctrl: ctrl}
	mock.recorder = &MockRegisterStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegisterStorage) EXPECT() *MockRegisterStorageMockRecorder {
	return m.recorder
}

// AddCredentials mocks base method.
func (m *MockRegisterStorage) AddCredentials(arg0 context.Context, arg1 entities.AuthCredentials) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCredentials", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCredentials indicates an expected call of AddCredentials.
func (mr *MockRegisterStorageMockRecorder) AddCredentials(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCredentials", reflect.TypeOf((*MockRegisterStorage)(nil).AddCredentials), arg0, arg1)
}

// FindUser mocks base method.
func (m *MockRegisterStorage) FindUser(arg0 context.Context, arg1 entities.Login) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// FindUser indicates an expected call of FindUser.
func (mr *MockRegisterStorageMockRecorder) FindUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*MockRegisterStorage)(nil).FindUser), arg0, arg1)
}