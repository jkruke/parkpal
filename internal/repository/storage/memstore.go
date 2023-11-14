package storage

import (
	"context"
	"parkpal-web-server/internal/entity"
)

var parkingLotList = []*entity.ParkingLot{
	{
		ID:   1,
		Name: "Hai Ba Trung",
	},
	{
		ID:   2,
		Name: "Dong Da",
	},
}

type memStore struct {
	db []*entity.ParkingLot
}

func NewMemStore() *memStore {
	return &memStore{db: make([]*entity.ParkingLot, 1000)}
}

func NewTestMemStore() *memStore {
	return &memStore{db: parkingLotList}
}

func (ms *memStore) GetParkingLotByID(ctx context.Context, id int) (*entity.ParkingLot, error) {
	i := ms.findIndexByParkingLotByID(id)
	if id == -1 {
		return nil, entity.ErrParkingLotNotFound
	}

	return ms.db[i], nil
}

func (ms *memStore) findIndexByParkingLotByID(id int) int {
	for i, p := range ms.db {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (ms *memStore) UpdateParkingLotByID(ctx context.Context, pl entity.ParkingLot) (*entity.ParkingLot, error) {
	i := ms.findIndexByParkingLotByID(pl.ID)
	if pl.ID == -1 {
		return nil, entity.ErrParkingLotNotFound
	}

	ms.db[i].BikeCount = pl.BikeCount
	ms.db[i].CongestionRate = pl.CongestionRate

	return ms.db[i], nil
}
