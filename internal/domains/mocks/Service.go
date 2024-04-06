// Code generated by mockery v3.0.0-alpha.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	shema "eCar/internal/shema"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// AddCar provides a mock function with given fields: ctx, cars
func (_m *Service) AddCar(ctx context.Context, cars []shema.Car) error {
	ret := _m.Called(ctx, cars)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []shema.Car) error); ok {
		r0 = rf(ctx, cars)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCar provides a mock function with given fields: ctx, id
func (_m *Service) DeleteCar(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetData provides a mock function with given fields: ctx, regNum, mark, model, year, ownerName, ownerSurname, ownerPatronymic, page, limit
func (_m *Service) GetData(ctx context.Context, regNum string, mark string, model string, year int, ownerName string, ownerSurname string, ownerPatronymic string, page int, limit int) ([]shema.Car, error) {
	ret := _m.Called(ctx, regNum, mark, model, year, ownerName, ownerSurname, ownerPatronymic, page, limit)

	var r0 []shema.Car
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, int, string, string, string, int, int) []shema.Car); ok {
		r0 = rf(ctx, regNum, mark, model, year, ownerName, ownerSurname, ownerPatronymic, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]shema.Car)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, int, string, string, string, int, int) error); ok {
		r1 = rf(ctx, regNum, mark, model, year, ownerName, ownerSurname, ownerPatronymic, page, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCar provides a mock function with given fields: ctx, id, filter
func (_m *Service) UpdateCar(ctx context.Context, id int, filter shema.Filter) error {
	ret := _m.Called(ctx, id, filter)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, shema.Filter) error); ok {
		r0 = rf(ctx, id, filter)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewService interface {
	mock.TestingT
	Cleanup(func())
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewService(t mockConstructorTestingTNewService) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}