// Code generated by MockGen. DO NOT EDIT.
// Source: ../auth/client.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	auth "github.com/reecerussell/open-social/client/auth"
	reflect "reflect"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockClient) GenerateToken(in *auth.GenerateTokenRequest) (*auth.GenerateTokenResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", in)
	ret0, _ := ret[0].(*auth.GenerateTokenResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockClientMockRecorder) GenerateToken(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockClient)(nil).GenerateToken), in)
}
