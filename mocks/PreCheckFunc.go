// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	types "github.com/tendermint/tendermint/types"
)

// PreCheckFunc is an autogenerated mock type for the PreCheckFunc type
type PreCheckFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *PreCheckFunc) Execute(_a0 types.Tx) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Tx) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPreCheckFunc interface {
	mock.TestingT
	Cleanup(func())
}

// NewPreCheckFunc creates a new instance of PreCheckFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPreCheckFunc(t mockConstructorTestingTNewPreCheckFunc) *PreCheckFunc {
	mock := &PreCheckFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
