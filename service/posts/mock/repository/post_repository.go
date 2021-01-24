// Code generated by MockGen. DO NOT EDIT.
// Source: ../repository/post_repository.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	model "github.com/reecerussell/open-social/service/posts/model"
	reflect "reflect"
)

// MockPostRepository is a mock of PostRepository interface.
type MockPostRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPostRepositoryMockRecorder
}

// MockPostRepositoryMockRecorder is the mock recorder for MockPostRepository.
type MockPostRepositoryMockRecorder struct {
	mock *MockPostRepository
}

// NewMockPostRepository creates a new mock instance.
func NewMockPostRepository(ctrl *gomock.Controller) *MockPostRepository {
	mock := &MockPostRepository{ctrl: ctrl}
	mock.recorder = &MockPostRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostRepository) EXPECT() *MockPostRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPostRepository) Create(ctx context.Context, p *model.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, p)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockPostRepositoryMockRecorder) Create(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPostRepository)(nil).Create), ctx, p)
}
