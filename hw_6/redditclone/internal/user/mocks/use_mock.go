// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/user/usecase.go

// Package mock_user is a generated GoMock package.
package mock_user

import (
	errors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	models "lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserUsecase is a mock of UserUsecase interface.
type MockUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseMockRecorder
}

// MockUserUsecaseMockRecorder is the mock recorder for MockUserUsecase.
type MockUserUsecaseMockRecorder struct {
	mock *MockUserUsecase
}

// NewMockUserUsecase creates a new mock instance.
func NewMockUserUsecase(ctrl *gomock.Controller) *MockUserUsecase {
	mock := &MockUserUsecase{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUsecase) EXPECT() *MockUserUsecaseMockRecorder {
	return m.recorder
}

// LoginUser mocks base method.
func (m *MockUserUsecase) LoginUser(user *models.User) (uint64, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", user)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockUserUsecaseMockRecorder) LoginUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockUserUsecase)(nil).LoginUser), user)
}

// RegiserUser mocks base method.
func (m *MockUserUsecase) RegiserUser(user *models.User) (uint64, *errors.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegiserUser", user)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(*errors.Error)
	return ret0, ret1
}

// RegiserUser indicates an expected call of RegiserUser.
func (mr *MockUserUsecaseMockRecorder) RegiserUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegiserUser", reflect.TypeOf((*MockUserUsecase)(nil).RegiserUser), user)
}
