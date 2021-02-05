// Code generated by MockGen. DO NOT EDIT.
// Source: ../repository/media_repository.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	model "github.com/reecerussell/open-social/service/media/model"
	reflect "reflect"
)

// MockMediaRepository is a mock of MediaRepository interface.
type MockMediaRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMediaRepositoryMockRecorder
}

// MockMediaRepositoryMockRecorder is the mock recorder for MockMediaRepository.
type MockMediaRepositoryMockRecorder struct {
	mock *MockMediaRepository
}

// NewMockMediaRepository creates a new mock instance.
func NewMockMediaRepository(ctrl *gomock.Controller) *MockMediaRepository {
	mock := &MockMediaRepository{ctrl: ctrl}
	mock.recorder = &MockMediaRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMediaRepository) EXPECT() *MockMediaRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m_2 *MockMediaRepository) Create(ctx context.Context, m *model.Media) (func(bool), error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Create", ctx, m)
	ret0, _ := ret[0].(func(bool))
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockMediaRepositoryMockRecorder) Create(ctx, m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMediaRepository)(nil).Create), ctx, m)
}

// GetContentType mocks base method.
func (m *MockMediaRepository) GetContentType(ctx context.Context, referenceID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContentType", ctx, referenceID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContentType indicates an expected call of GetContentType.
func (mr *MockMediaRepositoryMockRecorder) GetContentType(ctx, referenceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContentType", reflect.TypeOf((*MockMediaRepository)(nil).GetContentType), ctx, referenceID)
}
