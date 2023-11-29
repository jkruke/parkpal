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

	api.l.Info("Successfully get the parking lot")

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

	api.l.Info("Successfully get the parking lot")
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
