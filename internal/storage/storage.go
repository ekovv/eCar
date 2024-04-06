package storage

import (
	"context"
	"database/sql"
	"eCar/config"
	"eCar/internal/shema"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"strings"
)

type DBStorage struct {
	conn *sql.DB
}

func NewDBStorage(config config.Config) (*DBStorage, error) {
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db %w", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate driver, %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"car", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to do migrate %w", err)
	}
	s := &DBStorage{
		conn: db,
	}

	return s, s.CheckConnection()
}

func (s *DBStorage) CheckConnection() error {
	if err := s.conn.Ping(); err != nil {
		return fmt.Errorf("failed to connect to db %w", err)
	}
	return nil
}

func (s *DBStorage) SaveCars(ctx context.Context, cars []shema.Car) error {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO Cars (regNum, mark, model, year, owner_name, 
                  owner_surname, owner_patronymic) VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, car := range cars {
		_, err := stmt.ExecContext(ctx, car.RegNum, car.Mark, car.Model, car.Year, car.Owner.Name, car.Owner.Surname, car.Owner.Patronymic)
		if err != nil {
			err := tx.Rollback()

			if err != nil {
				return fmt.Errorf("failed to rollback")
			}
			return fmt.Errorf("can't save car in db: %w", err)
		}
	}

	return tx.Commit()
}

func (s *DBStorage) GetCars(ctx context.Context, regNum, mark, model string, year int, ownerName, ownerSurname,
	ownerPatronymic string, page, limit int) ([]shema.Car, error) {

	offset := (page - 1) * limit

	rows, err := s.conn.QueryContext(ctx, `SELECT regNum, mark, model, year, owner_name, owner_surname, owner_patronymic FROM Cars WHERE 
		(regNum = COALESCE($1, regNum)) OR 
		(mark = COALESCE($2, mark)) OR 
		(model = COALESCE($3, model)) OR
		(year = COALESCE($4, year)) OR
		(owner_name = COALESCE($5, owner_name)) OR
		(owner_surname = COALESCE($6, owner_surname)) OR
		(owner_patronymic = COALESCE($7, owner_patronymic)) 
		LIMIT $8 OFFSET $9`,
		regNum, mark, model, year, ownerName, ownerSurname, ownerPatronymic, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	defer rows.Close()

	var cars []shema.Car
	for rows.Next() {
		var car shema.Car

		err := rows.Scan(&car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		cars = append(cars, car)
	}

	return cars, nil
}

func (s *DBStorage) DeleteCar(ctx context.Context, id int) error {
	_, err := s.conn.ExecContext(ctx, "DELETE FROM Cars WHERE ID = $1", id)
	return err
}

func (s *DBStorage) UpdateCar(ctx context.Context, id int, filter shema.Filter) error {
	fields := make([]string, 0)
	values := make([]interface{}, 0)
	count := 1

	if filter.RegNum != "" {
		fields = append(fields, fmt.Sprintf("regNum = $%d", count))
		values = append(values, filter.RegNum)
		count++
	}
	if filter.Mark != "" {
		fields = append(fields, fmt.Sprintf("mark = $%d", count))
		values = append(values, filter.Mark)
		count++
	}
	if filter.Model != "" {
		fields = append(fields, fmt.Sprintf("model = $%d", count))
		values = append(values, filter.Model)
		count++
	}
	if filter.Year != 0 {
		fields = append(fields, fmt.Sprintf("year = $%d", count))
		values = append(values, filter.Year)
		count++
	}
	if filter.OwnerName != "" {
		fields = append(fields, fmt.Sprintf("owner_name = $%d", count))
		values = append(values, filter.OwnerName)
		count++
	}
	if filter.OwnerSurname != "" {
		fields = append(fields, fmt.Sprintf("owner_surname = $%d", count))
		values = append(values, filter.OwnerSurname)
		count++
	}
	if filter.OwnerPatronymic != "" {
		fields = append(fields, fmt.Sprintf("owner_patronymic = $%d", count))
		values = append(values, filter.OwnerPatronymic)
		count++
	}

	if len(fields) == 0 {
		return fmt.Errorf("field is empty")
	}

	query := fmt.Sprintf("UPDATE Cars SET %s WHERE ID = $%d", strings.Join(fields, ", "), count)
	values = append(values, id)

	_, err := s.conn.ExecContext(ctx, query, values...)
	return err
}

func (s *DBStorage) ShutDown() error {
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("error closing db: %w", err)
	}

	return nil

}
