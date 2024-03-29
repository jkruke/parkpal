package api

import (
	"net/http"
	"parkpal-web-server/internal/business"
	"parkpal-web-server/internal/entity"
	"parkpal-web-server/pkg"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type api struct {
	business.Business
	l hclog.Logger
}

type GenericError struct {
	Message string
}

func NewAPI(b business.Business, l hclog.Logger) *api {
	return &api{b, l}
}

func getParkingLotID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}

func (api *api) SearchBike(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	queryParams := r.URL.Query().Get("license_plate")
	prod, err := api.Business.SearchBike(r.Context(), &business.SearchBikeRequest{LicensePlate: queryParams})

	switch err {
	case nil:

	case entity.ErrParkingLotNotFound:
		api.l.Error("Unable to fetch bike", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		api.l.Error("Unable to fetching bike", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = pkg.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		api.l.Error("Unable to serializing bike", err)
	}

	api.l.Info("Successfully search the bike")

}

func (api *api) SearchParkingLot(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	queryParams := r.URL.Query().Get("name")
	prod, err := api.Business.SearchParkingLot(r.Context(), queryParams)

	switch err {
	case nil:

	case entity.ErrParkingLotNotFound:
		api.l.Error("Unable to fetch parking lot", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		api.l.Error("Unable to fetching parking lot", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = pkg.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		api.l.Error("Unable to serializing parking lot", err)
	}

	api.l.Info("Successfully search the parking lot")

}

func (api *api) GetAllParkingLots(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	prod, err := api.Business.GetAllParkingLots(r.Context())

	switch err {
	case nil:

	case entity.ErrParkingLotNotFound:
		api.l.Error("Unable to fetch parking lot", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		api.l.Error("Unable to fetching parking lot", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = pkg.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		api.l.Error("Unable to serializing parking lot", err)
	}

	api.l.Info("Successfully get the parking lots")
}

func (api *api) GetParkingLot(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getParkingLotID(r)

	prod, err := api.Business.GetParkingLot(
		r.Context(),
		&business.GetParkingLotRequest{
			ID: id,
		})

	switch err {
	case nil:

	case entity.ErrParkingLotNotFound:
		api.l.Error("Unable to fetch parking lot", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		api.l.Error("Unable to fetching parking lot", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = pkg.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		api.l.Error("Unable to serializing parking lot", err)
	}

	api.l.Info("Successfully get the parking lot")
}

func (api *api) UpdateParkingLot(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getParkingLotID(r)

	var updateRequest business.UpdateParkingLotRequest
	err := pkg.FromJSON(&updateRequest, r.Body)
	if err != nil {
		api.l.Error("Error parsing JSON payload", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		pkg.ToJSON(&GenericError{Message: "Invalid JSON payload"}, rw)
		return
	}

	// Add ID to the update request
	updateRequest.ID = id

	prod, err := api.Business.UpdateParkingLot(r.Context(), &updateRequest)

	switch err {
	case nil:
		// Handle success
		break

	case entity.ErrParkingLotNotFound:
		api.l.Error("Unable to fetch parking lot", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		api.l.Error("Unable to fetch parking lot", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = pkg.ToJSON(prod, rw)
	if err != nil {
		api.l.Error("Unable to serialize parking lot", "error", err)
	}

	api.l.Info("Successfully updated the parking lot")
}

func (api *api) AddParkingLot(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	var addRequest business.AddParkingLotRequest
	err := pkg.FromJSON(&addRequest, r.Body)
	if err != nil {
		api.l.Error("Error parsing JSON payload", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		pkg.ToJSON(&GenericError{Message: "Invalid JSON payload"}, rw)
		return
	}

	// Add ID to the update request

	prod, err := api.Business.AddParkingLot(r.Context(), &addRequest)

	switch err {
	case nil:
		// Handle success
		break

	case entity.ErrParkingLotNotFound:
		api.l.Error("Unable to fetch parking lot", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		api.l.Error("Unable to fetch parking lot", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = pkg.ToJSON(prod, rw)
	if err != nil {
		api.l.Error("Unable to serialize parking lot", "error", err)
	}

	api.l.Info("Successfully added the parking lot")
}

func (api *api) DeleteParkingLot(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getParkingLotID(r)

	var delRequest business.DeleteParkingLotRequest

	// Add ID to the update request
	delRequest.ID = id

	prod, err := api.Business.DeleteParkingLot(r.Context(), &delRequest)

	switch err {
	case nil:
		// Handle success
		break

	case entity.ErrParkingLotNotFound:
		api.l.Error("Unable to fetch parking lot", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		api.l.Error("Unable to fetch parking lot", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		pkg.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = pkg.ToJSON(prod, rw)
	if err != nil {
		api.l.Error("Unable to serialize parking lot", "error", err)
	}

	api.l.Info("Successfully deleted the parking lot")
}
