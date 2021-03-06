// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package repository is a generated GoMock package.
package repository

import (
	gomock "github.com/golang/mock/gomock"
	tour "github.com/truewebber/tourtoster/tour"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Features mocks base method
func (m *MockRepository) Features() ([]tour.Feature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Features")
	ret0, _ := ret[0].([]tour.Feature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Features indicates an expected call of Features
func (mr *MockRepositoryMockRecorder) Features() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Features", reflect.TypeOf((*MockRepository)(nil).Features))
}

// List mocks base method
func (m *MockRepository) List(arg0 *tour.Order, arg1 ...tour.Filter) ([]tour.Tour, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].([]tour.Tour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockRepositoryMockRecorder) List(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepository)(nil).List), varargs...)
}

// Tour mocks base method
func (m *MockRepository) Tour(ID int64) (*tour.Tour, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tour", ID)
	ret0, _ := ret[0].(*tour.Tour)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Tour indicates an expected call of Tour
func (mr *MockRepositoryMockRecorder) Tour(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tour", reflect.TypeOf((*MockRepository)(nil).Tour), ID)
}

// Save mocks base method
func (m *MockRepository) Save(t *tour.Tour) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", t)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockRepositoryMockRecorder) Save(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepository)(nil).Save), t)
}

// Delete mocks base method
func (m *MockRepository) Delete(t *tour.Tour) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", t)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRepositoryMockRecorder) Delete(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), t)
}
