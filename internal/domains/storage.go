package domains

import (
	"context"
	"eCar/internal/shema"
)

type Storage interface {
	SaveCars(ctx context.Context, cars []shema.Car) error
}
