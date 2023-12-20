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
	query := "SELECT * FROM parking_lots"

	// Query to retrieve multiple rows
	rows, err := ms.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err // Return error to the caller
	}
	defer rows.Close() // Close the rows when done to free up resources

	var parkingLots []*entity.ParkingLot // Dynamic slice capacity
	for rows.Next() {
		parkingLot := &entity.ParkingLot{} // Create a new instance for each row
		if err := rows.Scan(&parkingLot.ID, &parkingLot.Name, &parkingLot.Latitude, &parkingLot.Longitude, &parkingLot.BikeCount, &parkingLot.TotalSpace, &parkingLot.CongestionRate); err != nil {
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
	query := "SELECT * FROM parking_lots WHERE name = $1"
	err := ms.db.QueryRowContext(ctx, query, name).Scan(&parkingLot.ID, &parkingLot.Name, &parkingLot.Latitude, &parkingLot.Longitude, &parkingLot.BikeCount, &parkingLot.TotalSpace, &parkingLot.CongestionRate)
	if err != nil {
		return nil, err // Return the error to the caller
	}

	return &parkingLot, nil
}

func (ms *memStore) GetParkingLotByID(ctx context.Context, id int) (*entity.ParkingLot, error) {
	parkingLot := entity.ParkingLot{}
	query := "SELECT * FROM parking_lots WHERE id = $1"
	err := ms.db.QueryRowContext(ctx, query, id).Scan(&parkingLot.ID, &parkingLot.Name, &parkingLot.Latitude, &parkingLot.Longitude, &parkingLot.BikeCount, &parkingLot.TotalSpace, &parkingLot.CongestionRate)
	if err != nil {
		return nil, err // Return the error to the caller
	}

	return &parkingLot, nil
}

func (ms *memStore) UpdateParkingLotByID(ctx context.Context, pl entity.ParkingLot) (*entity.ParkingLot, error) {
	// var parkingLot = entity.ParkingLot{}
	// query := "UPDATE parking_lots SET  = $1 WHERE id = $2 RETURNING id, username, email"
	// err := r.db.QueryRowContext(ctx, query, newUsername, id).Scan(&user.ID, &user.Username, &user.Email)
	// if err != nil {
	// 	return &User{}, err
	// }
	//
	// return &user, nil
	return &entity.ParkingLot{}, nil
}
