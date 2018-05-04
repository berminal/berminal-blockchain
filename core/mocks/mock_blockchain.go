// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/iost-official/prototype/core/block (interfaces: Chain)

// Package core_mock is a generated GoMock package.
package core_mock

import (
	gomock "github.com/golang/mock/gomock"
	block "github.com/iost-official/prototype/core/block"
	reflect "reflect"
)

// MockChain is a mock of Chain interface
type MockChain struct {
	ctrl     *gomock.Controller
	recorder *MockChainMockRecorder
}

// MockChainMockRecorder is the mock recorder for MockChain
type MockChainMockRecorder struct {
	mock *MockChain
}

// NewMockChain creates a new mock instance
func NewMockChain(ctrl *gomock.Controller) *MockChain {
	mock := &MockChain{ctrl: ctrl}
	mock.recorder = &MockChainMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockChain) EXPECT() *MockChainMockRecorder {
	return m.recorder
}

// Iterator mocks base method
func (m *MockChain) Iterator() block.ChainIterator {
	ret := m.ctrl.Call(m, "Iterator")
	ret0, _ := ret[0].(block.ChainIterator)
	return ret0
}

// Iterator indicates an expected call of Iterator
func (mr *MockChainMockRecorder) Iterator() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Iterator", reflect.TypeOf((*MockChain)(nil).Iterator))
}

// Length mocks base method
func (m *MockChain) Length() int {
	ret := m.ctrl.Call(m, "Length")
	ret0, _ := ret[0].(int)
	return ret0
}

// Length indicates an expected call of Length
func (mr *MockChainMockRecorder) Length() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Length", reflect.TypeOf((*MockChain)(nil).Length))
}

// Push mocks base method
func (m *MockChain) Push(arg0 *block.Block) error {
	ret := m.ctrl.Call(m, "Push", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Push indicates an expected call of Push
func (mr *MockChainMockRecorder) Push(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockChain)(nil).Push), arg0)
}

// Top mocks base method
func (m *MockChain) Top() *block.Block {
	ret := m.ctrl.Call(m, "Top")
	ret0, _ := ret[0].(*block.Block)
	return ret0
}

// Top indicates an expected call of Top
func (mr *MockChainMockRecorder) Top() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Top", reflect.TypeOf((*MockChain)(nil).Top))
}
