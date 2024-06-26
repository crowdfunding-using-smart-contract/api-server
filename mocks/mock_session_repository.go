// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/datasource/repository/session_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	entity "fund-o/api-server/internal/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockSessionRepository is a mock of SessionRepository interface.
type MockSessionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepositoryMockRecorder
}

// MockSessionRepositoryMockRecorder is the mock recorder for MockSessionRepository.
type MockSessionRepositoryMockRecorder struct {
	mock *MockSessionRepository
}

// NewMockSessionRepository creates a new mock instance.
func NewMockSessionRepository(ctrl *gomock.Controller) *MockSessionRepository {
	mock := &MockSessionRepository{ctrl: ctrl}
	mock.recorder = &MockSessionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionRepository) EXPECT() *MockSessionRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSessionRepository) Create(session *entity.Session) (*entity.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", session)
	ret0, _ := ret[0].(*entity.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSessionRepositoryMockRecorder) Create(session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionRepository)(nil).Create), session)
}

// FindByID mocks base method.
func (m *MockSessionRepository) FindByID(id uuid.UUID) (*entity.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(*entity.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockSessionRepositoryMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockSessionRepository)(nil).FindByID), id)
}
