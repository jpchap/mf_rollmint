// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	dalc "github.com/celestiaorg/rollmint/proto/pb/dalc"
	mock "github.com/stretchr/testify/mock"
)

// DALCServiceServer is an autogenerated mock type for the DALCServiceServer type
type DALCServiceServer struct {
	mock.Mock
}

// CheckBlockAvailability provides a mock function with given fields: _a0, _a1
func (_m *DALCServiceServer) CheckBlockAvailability(_a0 context.Context, _a1 *dalc.CheckBlockAvailabilityRequest) (*dalc.CheckBlockAvailabilityResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *dalc.CheckBlockAvailabilityResponse
	if rf, ok := ret.Get(0).(func(context.Context, *dalc.CheckBlockAvailabilityRequest) *dalc.CheckBlockAvailabilityResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dalc.CheckBlockAvailabilityResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dalc.CheckBlockAvailabilityRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveBlocks provides a mock function with given fields: _a0, _a1
func (_m *DALCServiceServer) RetrieveBlocks(_a0 context.Context, _a1 *dalc.RetrieveBlocksRequest) (*dalc.RetrieveBlocksResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *dalc.RetrieveBlocksResponse
	if rf, ok := ret.Get(0).(func(context.Context, *dalc.RetrieveBlocksRequest) *dalc.RetrieveBlocksResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dalc.RetrieveBlocksResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dalc.RetrieveBlocksRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubmitBlock provides a mock function with given fields: _a0, _a1
func (_m *DALCServiceServer) SubmitBlock(_a0 context.Context, _a1 *dalc.SubmitBlockRequest) (*dalc.SubmitBlockResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *dalc.SubmitBlockResponse
	if rf, ok := ret.Get(0).(func(context.Context, *dalc.SubmitBlockRequest) *dalc.SubmitBlockResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dalc.SubmitBlockResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dalc.SubmitBlockRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewDALCServiceServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewDALCServiceServer creates a new instance of DALCServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDALCServiceServer(t mockConstructorTestingTNewDALCServiceServer) *DALCServiceServer {
	mock := &DALCServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}