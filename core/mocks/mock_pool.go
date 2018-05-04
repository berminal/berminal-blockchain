// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/iost-official/prototype/core/state (interfaces: Pool)

// Package core_mock is a generated GoMock package.
package core_mock

import (
	gomock "github.com/golang/mock/gomock"
	state "github.com/iost-official/prototype/core/state"
	reflect "reflect"
)

// MockPool is a mock of Pool interface
type MockPool struct {
	ctrl     *gomock.Controller
	recorder *MockPoolMockRecorder
}

// MockPoolMockRecorder is the mock recorder for MockPool
type MockPoolMockRecorder struct {
	mock *MockPool
}

// NewMockPool creates a new mock instance
func NewMockPool(ctrl *gomock.Controller) *MockPool {
	mock := &MockPool{ctrl: ctrl}
	mock.recorder = &MockPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPool) EXPECT() *MockPoolMockRecorder {
	return m.recorder
}

// Copy mocks base method
func (m *MockPool) Copy() state.Pool {
	ret := m.ctrl.Call(m, "Copy")
	ret0, _ := ret[0].(state.Pool)
	return ret0
}

// Copy indicates an expected call of Copy
func (mr *MockPoolMockRecorder) Copy() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockPool)(nil).Copy))
}

// Delete mocks base method
func (m *MockPool) Delete(arg0 state.Key) error {
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockPoolMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPool)(nil).Delete), arg0)
}

// Flush mocks base method
func (m *MockPool) Flush() error {
	ret := m.ctrl.Call(m, "Flush")
	ret0, _ := ret[0].(error)
	return ret0
}

// Flush indicates an expected call of Flush
func (mr *MockPoolMockRecorder) Flush() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockPool)(nil).Flush))
}

// Get mocks base method
func (m *MockPool) Get(arg0 state.Key) (state.Value, error) {
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(state.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockPoolMockRecorder) Get(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPool)(nil).Get), arg0)
}

// GetHM mocks base method
func (m *MockPool) GetHM(arg0, arg1 state.Key) (state.Value, error) {
	ret := m.ctrl.Call(m, "GetHM", arg0, arg1)
	ret0, _ := ret[0].(state.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHM indicates an expected call of GetHM
func (mr *MockPoolMockRecorder) GetHM(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHM", reflect.TypeOf((*MockPool)(nil).GetHM), arg0, arg1)
}

// GetPatch mocks base method
func (m *MockPool) GetPatch() state.Patch {
	ret := m.ctrl.Call(m, "GetPatch")
	ret0, _ := ret[0].(state.Patch)
	return ret0
}

// GetPatch indicates an expected call of GetPatch
func (mr *MockPoolMockRecorder) GetPatch() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPatch", reflect.TypeOf((*MockPool)(nil).GetPatch))
}

// Has mocks base method
func (m *MockPool) Has(arg0 state.Key) (bool, error) {
	ret := m.ctrl.Call(m, "Has", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Has indicates an expected call of Has
func (mr *MockPoolMockRecorder) Has(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*MockPool)(nil).Has), arg0)
}

// Put mocks base method
func (m *MockPool) Put(arg0 state.Key, arg1 state.Value) error {
	ret := m.ctrl.Call(m, "Put", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put
func (mr *MockPoolMockRecorder) Put(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockPool)(nil).Put), arg0, arg1)
}

// PutHM mocks base method
func (m *MockPool) PutHM(arg0, arg1 state.Key, arg2 state.Value) error {
	ret := m.ctrl.Call(m, "PutHM", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutHM indicates an expected call of PutHM
func (mr *MockPoolMockRecorder) PutHM(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutHM", reflect.TypeOf((*MockPool)(nil).PutHM), arg0, arg1, arg2)
}
