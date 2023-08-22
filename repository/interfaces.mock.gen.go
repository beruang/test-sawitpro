// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

package repository

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return _m.recorder
}

// FindByPhone mocks base method
func (_m *MockRepositoryInterface) FindByPhone(ctx context.Context, input FindByPhoneInput) (FindByPhoneOutput, error) {
	ret := _m.ctrl.Call(_m, "FindByPhone", ctx, input)
	ret0, _ := ret[0].(FindByPhoneOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPhone indicates an expected call of FindByPhone
func (_mr *MockRepositoryInterfaceMockRecorder) FindByPhone(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "FindByPhone", reflect.TypeOf((*MockRepositoryInterface)(nil).FindByPhone), arg0, arg1)
}

// FindBySlug mocks base method
func (_m *MockRepositoryInterface) FindBySlug(ctx context.Context, input FindBySlugInput) (FindBySlugOutput, error) {
	ret := _m.ctrl.Call(_m, "FindBySlug", ctx, input)
	ret0, _ := ret[0].(FindBySlugOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBySlug indicates an expected call of FindBySlug
func (_mr *MockRepositoryInterfaceMockRecorder) FindBySlug(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "FindBySlug", reflect.TypeOf((*MockRepositoryInterface)(nil).FindBySlug), arg0, arg1)
}

// Store mocks base method
func (_m *MockRepositoryInterface) Store(ctx context.Context, input RegistrationInput) (RegistrationOutput, error) {
	ret := _m.ctrl.Call(_m, "Store", ctx, input)
	ret0, _ := ret[0].(RegistrationOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Store indicates an expected call of Store
func (_mr *MockRepositoryInterfaceMockRecorder) Store(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Store", reflect.TypeOf((*MockRepositoryInterface)(nil).Store), arg0, arg1)
}

// Put mocks base method
func (_m *MockRepositoryInterface) Put(ctx context.Context, input UpdateUserInput) error {
	ret := _m.ctrl.Call(_m, "Put", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put
func (_mr *MockRepositoryInterfaceMockRecorder) Put(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Put", reflect.TypeOf((*MockRepositoryInterface)(nil).Put), arg0, arg1)
}