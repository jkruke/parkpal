package storage

import (
	"context"
	"database/sql"
	"parkpal-web-server/internal/entity"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type memStore struct {
	db DBTX
}

func NewMemStore(db DBTX) *memStore {
	return &memStore{db: db}
}

func (ms *memStore) GetAllParkingLots(ctx context.Context) ([]*entity.ParkingLot, error) {
	query := "SELECT pl.*, COALESCE(count(lp.id), 0) FROM parking_lots pl LEFT JOIN license_plates lp ON lp.parkinglot_id = pl.id GROUP BY 1"

	// Query to retrieve multiple rows
	rows, err := ms.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err // Return error to the caller
	}
	defer rows.Close() // Close the rows when done to free up resources

	var parkingLots []*entity.ParkingLot // Dynamic slice capacity
	for rows.Next() {
		parkingLot := &entity.ParkingLot{} // Create a new instance for each row
		if err := rows.Scan(&parkingLot.ID, &parkingLot.Name, &parkingLot.Latitude, &parkingLot.Longitude, &parkingLot.CongestionRate, &parkingLot.TotalSpace, &parkingLot.BikeCount); err != nil {
			return nil, err // Return error to the caller
		}
		parkingLots = append(parkingLots, parkingLot)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err // Return error to the caller
	}

	return parkingLots, nil
}

func (ms *memStore) GetParkingLotByName(ctx context.Context, name string) (*entity.ParkingLot, error) {
	parkingLot := entity.ParkingLot{}

	query := "SELECT pl.id, pl.name, pl.latitude, pl.longitude, pl.totalSpace, pl.congestionRate, count(pl.id) FROM parking_lots pl LEFT JOIN license_plates lp ON lp.parkinglot_id = pl.id WHERE pl.name = $1 GROUP BY 1"

	err := ms.db.QueryRowContext(ctx, query, name).Scan(
		&parkingLot.ID,
		&parkingLot.Name,
		&parkingLot.Latitude,
		&parkingLot.Longitude,
		&parkingLot.TotalSpace,
		&parkingLot.CongestionRate,
		&parkingLot.BikeCount,
	)

	if err == sql.ErrNoRows {
		return nil, entity.ErrParkingLotNotFound
	} else if err != nil {
		return nil, err
	}

	return &parkingLot, nil
}

func (ms *memStore) GetParkingLotByID(ctx context.Context, id int) (*entity.ParkingLot, error) {
	parkingLot := entity.ParkingLot{}
	query := "SELECT pl.*, COALESCE(count(lp.id), 0) FROM parking_lots pl LEFT JOIN license_plates lp ON lp.parkinglot_id = pl.id WHERE pl.id = $1 GROUP BY 1"
	err := ms.db.QueryRowContext(ctx, query, id).Scan(&parkingLot.ID, &parkingLot.Name, &parkingLot.Latitude, &parkingLot.Longitude, &parkingLot.CongestionRate, &parkingLot.TotalSpace, &parkingLot.BikeCount)
	if err != nil {
		return nil, err // Return the error to the caller
	}

	return &parkingLot, nil
}

func (ms *memStore) UpdateParkingLotByID(ctx context.Context, pl entity.ParkingLot) (*entity.ParkingLot, error) {
	var parkingLot = entity.ParkingLot{}
	query := "UPDATE parking_lots SET name = $1, latitude = $2, longitude = $3, totalSpace = $4 WHERE id = $5 RETURNING id, name, latitude, longitude, totalSpace"
	err := ms.db.QueryRowContext(ctx, query, pl.Name, pl.Latitude, pl.Longitude, pl.TotalSpace, pl.ID).Scan(&parkingLot.ID, &parkingLot.Name, &parkingLot.Latitude, &parkingLot.Longitude, &parkingLot.TotalSpace)
	if err != nil {
		return &entity.ParkingLot{}, err
	}

	return &parkingLot, nil
}

func (ms *memStore) DeleteParkingLotByID(ctx context.Context, pl entity.ParkingLot) (*entity.ParkingLot, error) {
	var parkingLot = entity.ParkingLot{}
	query := "DELETE FROM parking_lots WHERE id = $1 RETURNING id, name, latitude, longitude, totalSpace"
	err := ms.db.QueryRowContext(ctx, query, pl.ID).Scan(&parkingLot.ID, &parkingLot.Name, &parkingLot.Latitude, &parkingLot.Longitude, &parkingLot.TotalSpace)
	if err != nil {
		return &entity.ParkingLot{}, err
	}

	return &parkingLot, nil
}

func (ms *memStore) AddParkingLot(ctx context.Context, pl entity.ParkingLot) (*entity.ParkingLot, error) {
	parkingLot := pl
	query := "INSERT INTO parking_lots (name, latitude, longitude, totalSpace, congestionRate) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, latitude, longitude, totalSpace"

	err := ms.db.QueryRowContext(ctx, query, pl.Name, pl.Latitude, pl.Longitude, pl.TotalSpace, 0).
		Scan(&parkingLot.ID, &parkingLot.Name, &parkingLot.Latitude, &parkingLot.Longitude, &parkingLot.TotalSpace)

	if err != nil {
		return &entity.ParkingLot{}, err
	}

	return &parkingLot, nil
}
