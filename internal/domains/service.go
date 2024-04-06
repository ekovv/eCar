package domains

import (
	"context"
	"eCar/internal/shema"
)

type Service interface {
	AddCar(ctx context.Context, cars []shema.Car) error
	GetData(ctx context.Context, regNum, mark, model string, year int, ownerName, ownerSurname,
		ownerPatronymic string, page, limit int) ([]shema.Car, error)
	DeleteCar(ctx context.Context, id int) error
	UpdateCar(ctx context.Context, id int, filter shema.Filter) error
}
