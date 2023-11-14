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

type UpdateParkingLotRequest entity.ParkingLot

type UpdateParkingLotResponse entity.ParkingLot

type Business interface {
	GetParkingLot(ctx context.Context, request *GetParkingLotRequest) (*GetParkingLotResponse, error)
	UpdateParkingLot(ctx context.Context, request *UpdateParkingLotRequest) (*UpdateParkingLotResponse, error)
}

type Repository interface {
	GetParkingLotByID(ctx context.Context, id int) (*entity.ParkingLot, error)
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

func (b *business) GetParkingLot(c context.Context, request *GetParkingLotRequest) (*GetParkingLotResponse, error) {
	ctx, cancel := context.WithTimeout(c, b.timeout)
	defer cancel()

	prod, err := b.repository.GetParkingLotByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &GetParkingLotResponse{prod.ID, prod.Name, prod.BikeCount, prod.CongestionRate}, nil
}

func (b *business) UpdateParkingLot(c context.Context, request *UpdateParkingLotRequest) (*UpdateParkingLotResponse, error) {
	ctx, cancel := context.WithTimeout(c, b.timeout)
	defer cancel()

	prod, err := b.repository.UpdateParkingLotByID(ctx, entity.ParkingLot(*request))
	if err != nil {
		return nil, err
	}

	return &UpdateParkingLotResponse{prod.ID, prod.Name, prod.BikeCount, prod.CongestionRate}, nil
}
