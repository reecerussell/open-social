// Code generated by MockGen. DO NOT EDIT.
// Source: ../provider/post_provider.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	dto "github.com/reecerussell/open-social/cmd/posts/dto"
	reflect "reflect"
)

// MockPostProvider is a mock of PostProvider interface.
type MockPostProvider struct {
	ctrl     *gomock.Controller
	recorder *MockPostProviderMockRecorder
}

// MockPostProviderMockRecorder is the mock recorder for MockPostProvider.
type MockPostProviderMockRecorder struct {
	mock *MockPostProvider
}

// NewMockPostProvider creates a new mock instance.
func NewMockPostProvider(ctrl *gomock.Controller) *MockPostProvider {
	mock := &MockPostProvider{ctrl: ctrl}
	mock.recorder = &MockPostProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostProvider) EXPECT() *MockPostProviderMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockPostProvider) Get(ctx context.Context, postReferenceID, userReferenceID string) (*dto.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, postReferenceID, userReferenceID)
	ret0, _ := ret[0].(*dto.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPostProviderMockRecorder) Get(ctx, postReferenceID, userReferenceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPostProvider)(nil).Get), ctx, postReferenceID, userReferenceID)
}

// GetProfileFeed mocks base method.
func (m *MockPostProvider) GetProfileFeed(ctx context.Context, username string, userReferenceID uuid.UUID) ([]*dto.FeedItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileFeed", ctx, username, userReferenceID)
	ret0, _ := ret[0].([]*dto.FeedItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileFeed indicates an expected call of GetProfileFeed.
func (mr *MockPostProviderMockRecorder) GetProfileFeed(ctx, username, userReferenceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileFeed", reflect.TypeOf((*MockPostProvider)(nil).GetProfileFeed), ctx, username, userReferenceID)
}
