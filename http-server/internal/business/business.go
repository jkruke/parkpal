package business

import (
	"context"
	"parkpal-web-server/internal/entity"
	"time"

	"github.com/hashicorp/go-hclog"
)

type GetParkingLotRequest struct {
	ID int `json:"id"`
}

type GetParkingLotResponse entity.ParkingLot

type GetAllParkingLotResponse []entity.ParkingLot

type SearchParkingLotResponse entity.ParkingLot

type UpdateParkingLotRequest entity.ParkingLot

type UpdateParkingLotResponse entity.ParkingLot

type Business interface {
	GetParkingLot(ctx context.Context, request *GetParkingLotRequest) (*GetParkingLotResponse, error)
	GetAllParkingLots(ctx context.Context) (*GetAllParkingLotResponse, error)
	SearchParkingLot(ctx context.Context, name string) (*SearchParkingLotResponse, error)
	UpdateParkingLot(ctx context.Context, request *UpdateParkingLotRequest) (*UpdateParkingLotResponse, error)
}

type Repository interface {
	GetParkingLotByID(ctx context.Context, id int) (*entity.ParkingLot, error)
	GetAllParkingLots(ctx context.Context) ([]*entity.ParkingLot, error)
	GetParkingLotByName(ctx context.Context, name string) (*entity.ParkingLot, error)
	UpdateParkingLotByID(ctx context.Context, pl entity.ParkingLot) (*entity.ParkingLot, error)
}

type business struct {
	repository Repository
	timeout    time.Duration
	l          hclog.Logger
}

func NewBusiness(repository Repository, timeout time.Duration, l hclog.Logger) *business {
	return &business{
		repository,
		timeout,
		l,
	}
}

func (b *business) SearchParkingLot(c context.Context, name string) (*SearchParkingLotResponse, error) {
	ctx, cancel := context.WithTimeout(c, b.timeout)
	defer cancel()

	prod, err := b.repository.GetParkingLotByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &SearchParkingLotResponse{prod.ID, prod.Name, prod.BikeCount, prod.CongestionRate, prod.TotalSpace, prod.Longitude, prod.Latitude}, nil
}

func (b *business) GetAllParkingLots(c context.Context) (*GetAllParkingLotResponse, error) {
	ctx, cancel := context.WithTimeout(c, b.timeout)
	defer cancel()

	prods, err := b.repository.GetAllParkingLots(ctx)
	if err != nil {
		return nil, err
	}

	// Create a new slice of entity.ParkingLot
	response := make(GetAllParkingLotResponse, len(prods))

	// Copy values from prodsPointers to prods
	for i, ptr := range prods {
		response[i] = *ptr
	}
	return &response, nil
}

func (b *business) GetParkingLot(c context.Context, request *GetParkingLotRequest) (*GetParkingLotResponse, error) {
	ctx, cancel := context.WithTimeout(c, b.timeout)
	defer cancel()

	prod, err := b.repository.GetParkingLotByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &GetParkingLotResponse{prod.ID, prod.Name, prod.BikeCount, prod.CongestionRate, prod.TotalSpace, prod.Longitude, prod.Latitude}, nil
}

func (b *business) UpdateParkingLot(c context.Context, request *UpdateParkingLotRequest) (*UpdateParkingLotResponse, error) {
	ctx, cancel := context.WithTimeout(c, b.timeout)
	defer cancel()

	prod, err := b.repository.UpdateParkingLotByID(ctx, entity.ParkingLot(*request))
	if err != nil {
		return nil, err
	}

	return &UpdateParkingLotResponse{prod.ID, prod.Name, prod.BikeCount, prod.CongestionRate, prod.TotalSpace, prod.Longitude, prod.Latitude}, nil
}
