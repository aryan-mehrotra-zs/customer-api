package customer

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/amehrotra/customer-api/errors"

	"github.com/gorilla/mux"

	"github.com/amehrotra/customer-api/models"
	"github.com/amehrotra/customer-api/services"
)

type handler struct {
	service services.Customer
}

func New(service services.Customer) handler {
	return handler{service: service}
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	customer, err := getData(w, r)
	if err != nil {
		setStatusCode(w, r, err, customer)

		return
	}

	customer, err = h.service.Create(customer)

	setStatusCode(w, r, err, customer)

}

func (h handler) GetByID(w http.ResponseWriter, r *http.Request) {

	id, err := getId(w, r)
	if err != nil {
		setStatusCode(w, r, err, nil)

		return
	}

	customer, err := h.service.Get(id)

	setStatusCode(w, r, err, customer)

}

func (h handler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var customer models.Customer

	err = json.Unmarshal(body, &customer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	param := mux.Vars(r)
	idParam := param["id"]

	customer.ID, err = strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	customer, err = h.service.Update(customer)
	switch err.(type) {
	case errors.EntityNotFound:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		writeResponse(w, customer)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.service.Delete(id)

	switch err.(type) {
	case errors.EntityNotFound:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
