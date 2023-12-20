package storage

//
// import (
// 	"context"
// 	"parkpal-web-server/internal/entity"
// )
//
// var parkingLotList = []*entity.ParkingLot{
// 	{
// 		ID:             1,
// 		Name:           "D9",
// 		BikeCount:      1000,
// 		CongestionRate: 80,
// 		TotalSpace:     1000,
// 		Latitude:       105.84414,
// 		Longitude:      21.00405,
// 	},
// 	{
// 		ID:             2,
// 		Name:           "C7",
// 		BikeCount:      200,
// 		CongestionRate: 20,
// 		TotalSpace:     2000,
// 		Latitude:       105.844203,
// 		Longitude:      21.005155,
// 	},
// 	{
// 		ID:             3,
// 		Name:           "D8",
// 		BikeCount:      323,
// 		CongestionRate: 60,
// 		TotalSpace:     500,
// 		Latitude:       105.842279,
// 		Longitude:      21.004127,
// 	},
// }
//
// type memStore struct {
// 	db []*entity.ParkingLot
// }
//
// func NewMemStore() *memStore {
// 	return &memStore{db: make([]*entity.ParkingLot, 1000)}
// }
//
// func NewTestMemStore() *memStore {
// 	return &memStore{db: parkingLotList}
// }
//
// func (ms *memStore) GetAllParkingLots(ctx context.Context) ([]*entity.ParkingLot, error) {
// 	return ms.db, nil
// }
//
// func (ms *memStore) GetParkingLotByName(ctx context.Context, name string) (*entity.ParkingLot, error) {
// 	i := ms.findIndexByParkingLotName(name)
// 	if i == -1 {
// 		return nil, entity.ErrParkingLotNotFound
// 	}
//
// 	return ms.db[i], nil
// }
//
// func (ms *memStore) GetParkingLotByID(ctx context.Context, id int) (*entity.ParkingLot, error) {
// 	i := ms.findIndexByParkingLotByID(id)
// 	if id == -1 {
// 		return nil, entity.ErrParkingLotNotFound
// 	}
//
// 	return ms.db[i], nil
// }
//
// func (ms *memStore) findIndexByParkingLotName(name string) int {
// 	for i, p := range ms.db {
// 		if p.Name == name {
// 			return i
// 		}
// 	}
//
// 	return -1
// }
//
// func (ms *memStore) findIndexByParkingLotByID(id int) int {
// 	for i, p := range ms.db {
// 		if p.ID == id {
// 			return i
// 		}
// 	}
//
// 	return -1
// }
//
// func (ms *memStore) UpdateParkingLotByID(ctx context.Context, pl entity.ParkingLot) (*entity.ParkingLot, error) {
// 	i := ms.findIndexByParkingLotByID(pl.ID)
// 	if pl.ID == -1 {
// 		return nil, entity.ErrParkingLotNotFound
// 	}
//
// 	ms.db[i].BikeCount = pl.BikeCount
// 	ms.db[i].CongestionRate = pl.CongestionRate
//
// 	return ms.db[i], nil
// }
