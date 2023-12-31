// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/sms/types.go
//
// Generated by this command:
//
//	mockgen -source=./internal/service/sms/types.go -package=smsmocks -destination=./internal/service/sms/mocks/sms.mock.go
//

// Package smsmocks is a generated GoMock package.
package smsmocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSMSService is a mock of SMSService interface.
type MockSMSService struct {
	ctrl     *gomock.Controller
	recorder *MockSMSServiceMockRecorder
}

// MockSMSServiceMockRecorder is the mock recorder for MockSMSService.
type MockSMSServiceMockRecorder struct {
	mock *MockSMSService
}

// NewMockSMSService creates a new mock instance.
func NewMockSMSService(ctrl *gomock.Controller) *MockSMSService {
	mock := &MockSMSService{ctrl: ctrl}
	mock.recorder = &MockSMSServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSMSService) EXPECT() *MockSMSServiceMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, tplId, args}
	for _, a := range numbers {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Send", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockSMSServiceMockRecorder) Send(ctx, tplId, args any, numbers ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, tplId, args}, numbers...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSMSService)(nil).Send), varargs...)
}
