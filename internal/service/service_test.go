package service

import (
	"context"
	"eCar/internal/constants"
	"eCar/internal/domains/mocks"
	"eCar/internal/shema"
	"errors"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

type storageMock func(c *mocks.Storage)

func TestService_AddCar(t *testing.T) {
	cars := []shema.Car{
		{
			RegNum: "DDDDDD",
			Mark:   "Honda",
			Model:  "Civic",
			Year:   2015,
			Owner: shema.People{
				Name:       "Petr",
				Surname:    "Petrov",
				Patronymic: "Petrovich",
			},
		},
	}
	tests := []struct {
		name        string
		storageMock storageMock
		wantErr     error
	}{
		{
			name: "OK1",
			storageMock: func(c *mocks.Storage) {
				c.Mock.On("SaveCars", mock.Anything, cars).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		{
			name: "BAD",
			storageMock: func(c *mocks.Storage) {
				c.Mock.On("SaveCars", mock.Anything, mock.Anything).Return(constants.ErrInvalidData).Times(1)
			},
			wantErr: constants.ErrInvalidData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)
			tt.storageMock(storage)
			logger, err := zap.NewProduction()

			service := Service{
				storage: storage,
				logger:  logger,
			}
			err = service.AddCar(context.Background(), cars)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetData(t *testing.T) {
	filter := shema.Filter{
		RegNum: "DDDDDD",
		Mark:   "Honda",
		Page:   1,
		Limit:  2,
	}

	tests := []struct {
		name        string
		storageMock storageMock
		wantErr     error
		want        []shema.Car
	}{
		{
			name: "OK1",
			storageMock: func(c *mocks.Storage) {
				c.Mock.On("GetCars", mock.Anything, filter.RegNum, filter.Mark, filter.Model, filter.Year, filter.OwnerName, filter.OwnerSurname,
					filter.OwnerPatronymic, filter.Page, filter.Limit).Return([]shema.Car{{
					RegNum: "DDDDDD",
					Mark:   "Honda",
					Model:  "Civic",
					Year:   2015,
					Owner: shema.People{
						Name:       "Petr",
						Surname:    "Petrov",
						Patronymic: "Petrovich",
					},
				}}, nil).Times(1)
			},
			wantErr: nil,
			want: []shema.Car{{
				RegNum: "DDDDDD",
				Mark:   "Honda",
				Model:  "Civic",
				Year:   2015,
				Owner: shema.People{
					Name:       "Petr",
					Surname:    "Petrov",
					Patronymic: "Petrovich",
				},
			}},
		},
		{
			name: "BAD",
			storageMock: func(c *mocks.Storage) {
				c.Mock.On("GetCars", mock.Anything, filter.RegNum, filter.Mark, filter.Model, filter.Year, filter.OwnerName, filter.OwnerSurname,
					filter.OwnerPatronymic, filter.Page, filter.Limit).Return(nil, constants.ErrInvalidData).Times(1)
			},
			wantErr: constants.ErrInvalidData,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)
			tt.storageMock(storage)
			logger, err := zap.NewProduction()

			service := Service{
				storage: storage,
				logger:  logger,
			}
			arr, err := service.GetData(context.Background(), filter.RegNum, filter.Mark, filter.Model, filter.Year, filter.OwnerName, filter.OwnerSurname,
				filter.OwnerPatronymic, filter.Page, filter.Limit)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(arr, tt.want) {
				t.Errorf("got %v, want %v", arr, tt.want)
			}
		})
	}
}

func TestService_DeleteCar(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		storageMock storageMock
		wantErr     error
	}{
		{
			name: "OK1",
			id:   2,
			storageMock: func(c *mocks.Storage) {
				c.Mock.On("DeleteCar", mock.Anything, 2).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		{
			name: "BAD",
			id:   12312,
			storageMock: func(c *mocks.Storage) {
				c.Mock.On("DeleteCar", mock.Anything, 12312).Return(constants.ErrInvalidData).Times(1)
			},
			wantErr: constants.ErrInvalidData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)
			tt.storageMock(storage)
			logger, err := zap.NewProduction()

			service := Service{
				storage: storage,
				logger:  logger,
			}
			err = service.DeleteCar(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_UpdateCar(t *testing.T) {
	filter := shema.Filter{
		RegNum: "DDDDDD",
		Mark:   "Honda",
	}

	tests := []struct {
		name        string
		id          int
		storageMock storageMock
		wantErr     error
	}{
		{
			name: "OK1",
			id:   2,
			storageMock: func(c *mocks.Storage) {
				c.Mock.On("UpdateCar", mock.Anything, 2, filter).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		{
			name: "BAD",
			id:   12312,
			storageMock: func(c *mocks.Storage) {
				c.Mock.On("UpdateCar", mock.Anything, 12312, filter).Return(constants.ErrInvalidData).Times(1)
			},
			wantErr: constants.ErrInvalidData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)
			tt.storageMock(storage)
			logger, err := zap.NewProduction()

			service := Service{
				storage: storage,
				logger:  logger,
			}
			err = service.UpdateCar(context.Background(), tt.id, filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
