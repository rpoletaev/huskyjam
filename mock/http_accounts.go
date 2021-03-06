// Code generated by MockGen. DO NOT EDIT.
// Source: ./http/accounts.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPassHashHelper is a mock of PassHashHelper interface.
type MockPassHashHelper struct {
	ctrl     *gomock.Controller
	recorder *MockPassHashHelperMockRecorder
}

// MockPassHashHelperMockRecorder is the mock recorder for MockPassHashHelper.
type MockPassHashHelperMockRecorder struct {
	mock *MockPassHashHelper
}

// NewMockPassHashHelper creates a new mock instance.
func NewMockPassHashHelper(ctrl *gomock.Controller) *MockPassHashHelper {
	mock := &MockPassHashHelper{ctrl: ctrl}
	mock.recorder = &MockPassHashHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPassHashHelper) EXPECT() *MockPassHashHelperMockRecorder {
	return m.recorder
}

// Hash mocks base method.
func (m *MockPassHashHelper) Hash(pass string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hash", pass)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hash indicates an expected call of Hash.
func (mr *MockPassHashHelperMockRecorder) Hash(pass interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hash", reflect.TypeOf((*MockPassHashHelper)(nil).Hash), pass)
}

// Check mocks base method.
func (m *MockPassHashHelper) Check(pass, hash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", pass, hash)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockPassHashHelperMockRecorder) Check(pass, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockPassHashHelper)(nil).Check), pass, hash)
}
