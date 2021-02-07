// Code generated by MockGen. DO NOT EDIT.
// Source: ../repository/like_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLikeRepository is a mock of LikeRepository interface.
type MockLikeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockLikeRepositoryMockRecorder
}

// MockLikeRepositoryMockRecorder is the mock recorder for MockLikeRepository.
type MockLikeRepositoryMockRecorder struct {
	mock *MockLikeRepository
}

// NewMockLikeRepository creates a new mock instance.
func NewMockLikeRepository(ctrl *gomock.Controller) *MockLikeRepository {
	mock := &MockLikeRepository{ctrl: ctrl}
	mock.recorder = &MockLikeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLikeRepository) EXPECT() *MockLikeRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockLikeRepository) Create(ctx context.Context, postID int, userReferenceID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, postID, userReferenceID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockLikeRepositoryMockRecorder) Create(ctx, postID, userReferenceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLikeRepository)(nil).Create), ctx, postID, userReferenceID)
}