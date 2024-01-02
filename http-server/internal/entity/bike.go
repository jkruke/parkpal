package entity

type Bike struct {
	ID           int        `json:"id"`
	LicensePlate string     `json:"plate"`
	ParkingLot   ParkingLot `json:"parking_lot"`
}
