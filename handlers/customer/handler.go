package customer

import (
	"encoding/json"
	"io"
	"log"
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
	w.Header().Set("Content-Type", "application/json")

	var customer models.Customer

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.Unmarshal(body, &customer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	customer, err = h.service.Create(customer)
	switch err.(type) {
	case errors.EntityAlreadyExists:
		w.WriteHeader(http.StatusOK)
	case errors.InvalidParam, errors.MissingParam:
		w.WriteHeader(http.StatusBadRequest)
	case nil:
		resp, err := json.Marshal(customer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusCreated)

		_, err = w.Write(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func writeResponse(w http.ResponseWriter, data interface{}) {
	res, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(res)
	if err != nil {
		log.Println(err)
	}
}

func (h handler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	data, err := h.service.Get(id)
	switch err.(type) {
	case errors.EntityNotFound:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		writeResponse(w, data)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
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
