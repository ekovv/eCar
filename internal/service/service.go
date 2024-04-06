package service

import (
	"context"
	"eCar/config"
	"eCar/internal/constants"
	"eCar/internal/domains"
	"eCar/internal/shema"
	"fmt"
	"go.uber.org/zap"
)

type Service struct {
	storage domains.Storage
	config  config.Config
	logger  *zap.Logger
}

func NewService(storage domains.Storage, config config.Config) *Service {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil
	}
	return &Service{storage: storage, config: config, logger: logger}
}

// AddCar add cars
func (s *Service) AddCar(ctx context.Context, cars []shema.Car) error {
	const op = "service.AddCar"

	err := s.storage.SaveCars(ctx, cars)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : failed to save cars: %v", op, err))
		return constants.ErrInvalidData
	}
	return nil
}

// GetData get cars from db
func (s *Service) GetData(ctx context.Context, regNum, mark, model string, year int, ownerName, ownerSurname,
	ownerPatronymic string, page, limit int) ([]shema.Car, error) {

	const op = "service.GetData"

	data, err := s.storage.GetCars(ctx, regNum, mark, model, year, ownerName, ownerSurname, ownerPatronymic, page, limit)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : failed to get data: %v", op, err))
		return nil, constants.ErrInvalidData
	}
	if len(data) == 0 {
		s.logger.Info(fmt.Sprintf("%s : data is empty: %v", op, err))
		return nil, constants.ErrInvalidData
	}
	return data, nil
}

// DeleteCar delete cars
func (s *Service) DeleteCar(ctx context.Context, id int) error {
	const op = "service.DeleteCar"

	err := s.storage.DeleteCar(ctx, id)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : failed to delete data: %v", op, err))
		return constants.ErrInvalidData
	}

	return nil
}

// UpdateCar update cars
func (s *Service) UpdateCar(ctx context.Context, id int, filter shema.Filter) error {
	const op = "service.UpdateCar"

	err := s.storage.UpdateCar(ctx, id, filter)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : failed to update data: %v", op, err))
		return constants.ErrInvalidData
	}

	return nil
}
