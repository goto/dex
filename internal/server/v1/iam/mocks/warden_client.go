// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	warden "github.com/goto/dex/warden"
)

// WardenClient is an autogenerated mock type for the WardenClient type
type WardenClient struct {
	mock.Mock
}

type WardenClient_Expecter struct {
	mock *mock.Mock
}

func (_m *WardenClient) EXPECT() *WardenClient_Expecter {
	return &WardenClient_Expecter{mock: &_m.Mock}
}

// ListUserTeams provides a mock function with given fields: ctx, req
func (_m *WardenClient) ListUserTeams(ctx context.Context, req warden.TeamListRequest) ([]warden.Team, error) {
	ret := _m.Called(ctx, req)

	var r0 []warden.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, warden.TeamListRequest) ([]warden.Team, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, warden.TeamListRequest) []warden.Team); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]warden.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, warden.TeamListRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WardenClient_ListUserTeams_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListUserTeams'
type WardenClient_ListUserTeams_Call struct {
	*mock.Call
}

// ListUserTeams is a helper method to define mock.On call
//   - ctx context.Context
//   - req warden.TeamListRequest
func (_e *WardenClient_Expecter) ListUserTeams(ctx interface{}, req interface{}) *WardenClient_ListUserTeams_Call {
	return &WardenClient_ListUserTeams_Call{Call: _e.mock.On("ListUserTeams", ctx, req)}
}

func (_c *WardenClient_ListUserTeams_Call) Run(run func(ctx context.Context, req warden.TeamListRequest)) *WardenClient_ListUserTeams_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(warden.TeamListRequest))
	})
	return _c
}

func (_c *WardenClient_ListUserTeams_Call) Return(_a0 []warden.Team, _a1 error) *WardenClient_ListUserTeams_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *WardenClient_ListUserTeams_Call) RunAndReturn(run func(context.Context, warden.TeamListRequest) ([]warden.Team, error)) *WardenClient_ListUserTeams_Call {
	_c.Call.Return(run)
	return _c
}

// TeamByUUID provides a mock function with given fields: ctx, req
func (_m *WardenClient) TeamByUUID(ctx context.Context, req warden.TeamByUUIDRequest) (*warden.Team, error) {
	ret := _m.Called(ctx, req)

	var r0 *warden.Team
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, warden.TeamByUUIDRequest) (*warden.Team, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, warden.TeamByUUIDRequest) *warden.Team); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*warden.Team)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, warden.TeamByUUIDRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WardenClient_TeamByUUID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TeamByUUID'
type WardenClient_TeamByUUID_Call struct {
	*mock.Call
}

// TeamByUUID is a helper method to define mock.On call
//   - ctx context.Context
//   - req warden.TeamByUUIDRequest
func (_e *WardenClient_Expecter) TeamByUUID(ctx interface{}, req interface{}) *WardenClient_TeamByUUID_Call {
	return &WardenClient_TeamByUUID_Call{Call: _e.mock.On("TeamByUUID", ctx, req)}
}

func (_c *WardenClient_TeamByUUID_Call) Run(run func(ctx context.Context, req warden.TeamByUUIDRequest)) *WardenClient_TeamByUUID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(warden.TeamByUUIDRequest))
	})
	return _c
}

func (_c *WardenClient_TeamByUUID_Call) Return(_a0 *warden.Team, _a1 error) *WardenClient_TeamByUUID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *WardenClient_TeamByUUID_Call) RunAndReturn(run func(context.Context, warden.TeamByUUIDRequest) (*warden.Team, error)) *WardenClient_TeamByUUID_Call {
	_c.Call.Return(run)
	return _c
}

// NewWardenClient creates a new instance of WardenClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWardenClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *WardenClient {
	mock := &WardenClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
