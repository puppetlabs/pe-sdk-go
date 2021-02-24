// Code generated by MockGen. DO NOT EDIT.
// Source: ../token.go

// Package mock_token is a generated GoMock package.
package mock_token

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockToken is a mock of Token interface
type MockToken struct {
	ctrl     *gomock.Controller
	recorder *MockTokenMockRecorder
}

// MockTokenMockRecorder is the mock recorder for MockToken
type MockTokenMockRecorder struct {
	mock *MockToken
}

// NewMockToken creates a new mock instance
func NewMockToken(ctrl *gomock.Controller) *MockToken {
	mock := &MockToken{ctrl: ctrl}
	mock.recorder = &MockTokenMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockToken) EXPECT() *MockTokenMockRecorder {
	return m.recorder
}

// Read mocks base method
func (m *MockToken) Read() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockTokenMockRecorder) Read() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockToken)(nil).Read))
}

// Write mocks base method
func (m *MockToken) Write(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write
func (mr *MockTokenMockRecorder) Write(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockToken)(nil).Write), arg0)
}

// Delete mocks base method
func (m *MockToken) Delete(arg0 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockTokenMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockToken)(nil).Delete), arg0)
}
