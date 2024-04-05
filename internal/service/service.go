package service

import (
	"context"
	"eCar/config"
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

func (s *Service) AddCar(ctx context.Context, cars []shema.Car) error {
	const op = "service.AddCar"

	err := s.storage.SaveCars(ctx, cars)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s : failed to save cars: %v", op, err))
		return err
	}
	return nil
}
