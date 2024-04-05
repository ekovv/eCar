package domains

import (
	"context"
	"eCar/internal/shema"
)

type Service interface {
	AddCar(ctx context.Context, cars []shema.Car) error
}
