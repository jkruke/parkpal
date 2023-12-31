package entity

import "errors"

var ErrNoMessagesFound = errors.New("no messages found")
var ErrParkingLotNotFound = errors.New("no parking lot found")

type ParkingLot struct {
	ID             int     `json:"id"` // Unique identifier for the product
	Name           string  `json:"name"`
	BikeCount      int     `json:"bikeCount"`
	CongestionRate float64 `json:"congestionRate"`
	TotalSpace     int     `json:"totalSpace"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
}
