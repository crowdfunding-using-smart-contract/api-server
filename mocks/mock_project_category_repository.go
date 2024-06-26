// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/datasource/repository/project_category_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	entity "fund-o/api-server/internal/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockProjectCategoryRepository is a mock of ProjectCategoryRepository interface.
type MockProjectCategoryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProjectCategoryRepositoryMockRecorder
}

// MockProjectCategoryRepositoryMockRecorder is the mock recorder for MockProjectCategoryRepository.
type MockProjectCategoryRepositoryMockRecorder struct {
	mock *MockProjectCategoryRepository
}

// NewMockProjectCategoryRepository creates a new mock instance.
func NewMockProjectCategoryRepository(ctrl *gomock.Controller) *MockProjectCategoryRepository {
	mock := &MockProjectCategoryRepository{ctrl: ctrl}
	mock.recorder = &MockProjectCategoryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectCategoryRepository) EXPECT() *MockProjectCategoryRepositoryMockRecorder {
	return m.recorder
}

// FindAll mocks base method.
func (m *MockProjectCategoryRepository) FindAll() ([]entity.ProjectCategory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll")
	ret0, _ := ret[0].([]entity.ProjectCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockProjectCategoryRepositoryMockRecorder) FindAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockProjectCategoryRepository)(nil).FindAll))
}
