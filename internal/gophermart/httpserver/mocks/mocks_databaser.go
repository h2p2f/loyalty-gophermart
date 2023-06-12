// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// DataBaser is an autogenerated mock type for the DataBaser type
type DataBaser struct {
	mock.Mock
}

// CheckUniqueOrder provides a mock function with given fields: order
func (_m *DataBaser) CheckUniqueOrder(order string) (string, bool) {
	ret := _m.Called(order)

	var r0 string
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (string, bool)); ok {
		return rf(order)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(order)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(order)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// FindPassByLogin provides a mock function with given fields: login
func (_m *DataBaser) FindPassByLogin(login string) (string, error) {
	ret := _m.Called(login)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(login)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(login)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllWithdraws provides a mock function with given fields: login
func (_m *DataBaser) GetAllWithdraws(login string) []byte {
	ret := _m.Called(login)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// GetBalance provides a mock function with given fields: login
func (_m *DataBaser) GetBalance(login string) (float64, error) {
	ret := _m.Called(login)

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (float64, error)); ok {
		return rf(login)
	}
	if rf, ok := ret.Get(0).(func(string) float64); ok {
		r0 = rf(login)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrdersByUser provides a mock function with given fields: login
func (_m *DataBaser) GetOrdersByUser(login string) ([]byte, error) {
	ret := _m.Called(login)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]byte, error)); ok {
		return rf(login)
	}
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSumOfAllWithdraws provides a mock function with given fields: login
func (_m *DataBaser) GetSumOfAllWithdraws(login string) float64 {
	ret := _m.Called(login)

	var r0 float64
	if rf, ok := ret.Get(0).(func(string) float64); ok {
		r0 = rf(login)
	} else {
		r0 = ret.Get(0).(float64)
	}

	return r0
}

// NewOrder provides a mock function with given fields: id, login, status, accrual, timeCreated
func (_m *DataBaser) NewOrder(id string, login string, status string, accrual float64, timeCreated time.Time) error {
	ret := _m.Called(id, login, status, accrual, timeCreated)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, float64, time.Time) error); ok {
		r0 = rf(id, login, status, accrual, timeCreated)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUser provides a mock function with given fields: login, password
func (_m *DataBaser) NewUser(login string, password string) error {
	ret := _m.Called(login, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(login, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWithdraw provides a mock function with given fields: login, order, amount, timeCreated
func (_m *DataBaser) NewWithdraw(login string, order string, amount float64, timeCreated time.Time) error {
	ret := _m.Called(login, order, amount, timeCreated)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, float64, time.Time) error); ok {
		r0 = rf(login, order, amount, timeCreated)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewDataBaser interface {
	mock.TestingT
	Cleanup(func())
}

// NewDataBaser creates a new instance of DataBaser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDataBaser(t mockConstructorTestingTNewDataBaser) *DataBaser {
	mock := &DataBaser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
