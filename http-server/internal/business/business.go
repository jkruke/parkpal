package business

import (
	"context"
	"errors"
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
type UpdateParkingLotResponse struct {
	ID         int     `json:"id"` // Unique identifier for the product
	Name       string  `json:"name"`
	TotalSpace int     `json:"totalSpace"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
}

type AddParkingLotRequest entity.ParkingLot
type AddParkingLotResponse struct {
	ID         int     `json:"id"` // Unique identifier for the product
	Name       string  `json:"name"`
	TotalSpace int     `json:"totalSpace"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
}

type DeleteParkingLotRequest entity.ParkingLot
type DeleteParkingLotResponse struct {
	ID         int     `json:"id"` // Unique identifier for the product
	Name       string  `json:"name"`
	TotalSpace int     `json:"totalSpace"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
}

type Business interface {
	GetParkingLot(ctx context.Context, request *GetParkingLotRequest) (*GetParkingLotResponse, error)
	GetAllParkingLots(ctx context.Context) (*GetAllParkingLotResponse, error)
	SearchParkingLot(ctx context.Context, name string) (*SearchParkingLotResponse, error)
	UpdateParkingLot(ctx context.Context, request *UpdateParkingLotRequest) (*UpdateParkingLotResponse, error)
	DeleteParkingLot(ctx context.Context, request *DeleteParkingLotRequest) (*DeleteParkingLotResponse, error)
	AddParkingLot(ctx context.Context, request *AddParkingLotRequest) (*AddParkingLotResponse, error)
}

type Repository interface {
	GetParkingLotByID(ctx context.Context, id int) (*entity.ParkingLot, error)
	GetAllParkingLots(ctx context.Context) ([]*entity.ParkingLot, error)
	GetParkingLotByName(ctx context.Context, name string) (*entity.ParkingLot, error)
	UpdateParkingLotByID(ctx context.Context, pl entity.ParkingLot) (*entity.ParkingLot, error)
	DeleteParkingLotByID(ctx context.Context, pl entity.ParkingLot) (*entity.ParkingLot, error)
	AddParkingLot(ctx context.Context, pl entity.ParkingLot) (*entity.ParkingLot, error)
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
		if ptr != nil {
			response[i] = *ptr
		} else {
			return nil, errors.New("received nil parking lot from repository")
		}
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

	return &UpdateParkingLotResponse{prod.ID, prod.Name, prod.TotalSpace, prod.Longitude, prod.Latitude}, nil
}

func (b *business) DeleteParkingLot(c context.Context, request *DeleteParkingLotRequest) (*DeleteParkingLotResponse, error) {
	ctx, cancel := context.WithTimeout(c, b.timeout)
	defer cancel()

	prod, err := b.repository.DeleteParkingLotByID(ctx, entity.ParkingLot(*request))
	if err != nil {
		return nil, err
	}

	return &DeleteParkingLotResponse{prod.ID, prod.Name, prod.TotalSpace, prod.Longitude, prod.Latitude}, nil
}

func (b *business) AddParkingLot(c context.Context, request *AddParkingLotRequest) (*AddParkingLotResponse, error) {
	ctx, cancel := context.WithTimeout(c, b.timeout)
	defer cancel()

	prod, err := b.repository.AddParkingLot(ctx, entity.ParkingLot(*request))
	if err != nil {
		return nil, err
	}

	return &AddParkingLotResponse{prod.ID, prod.Name, prod.TotalSpace, prod.Longitude, prod.Latitude}, nil
}
