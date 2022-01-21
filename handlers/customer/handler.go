package customer

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/amehrotra/customer-api/errors"
	"github.com/amehrotra/customer-api/models"
	"github.com/amehrotra/customer-api/services"
)

type handler struct {
	service services.Customer
}

// nolint:revive // handler should be unexported
func New(service services.Customer) handler {
	return handler{service: service}
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	customer, err := getCustomer(r)
	if err != nil {
		setStatusCode(w, r.Method, nil, err)

		return
	}

	res, err := h.service.Create(customer)

	setStatusCode(w, r.Method, res, err)
}

func (h handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		setStatusCode(w, r.Method, nil, err)

		return
	}

	customer, err := h.service.Get(id)

	setStatusCode(w, r.Method, customer, err)
}

func (h handler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	customer, err := getCustomer(r)
	if err != nil {
		setStatusCode(w, r.Method, nil, err)

		return
	}

	customer.ID, err = getID(r)
	if err != nil {
		setStatusCode(w, r.Method, nil, err)

		return
	}

	res, err := h.service.Update(customer)

	setStatusCode(w, r.Method, res, err)
}

func (h handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		setStatusCode(w, r.Method, nil, err)

		return
	}

	err = h.service.Delete(id)

	setStatusCode(w, r.Method, nil, err)
}

func getID(r *http.Request) (int, error) {
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return 0, errors.InvalidParam{Param: []string{"id"}}
	}

	return id, nil
}

func getCustomer(r *http.Request) (models.Customer, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return models.Customer{}, errors.Error("bind error")
	}

	var customer models.Customer

	err = json.Unmarshal(body, &customer)
	if err != nil {
		return models.Customer{}, errors.InvalidParam{Param: []string{"body"}}
	}

	return customer, nil
}

func setStatusCode(w http.ResponseWriter, method string, data interface{}, err error) {
	switch err.(type) {
	case errors.EntityAlreadyExists:
		w.WriteHeader(http.StatusOK)
	case errors.InvalidParam, errors.MissingParam:
		w.WriteHeader(http.StatusBadRequest)
	case errors.EntityNotFound:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		writeSuccessResponse(w, method, data)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func writeSuccessResponse(w http.ResponseWriter, method string, data interface{}) {
	switch method {
	case http.MethodPost:
		writeResponseBody(w, http.StatusCreated, data)
	case http.MethodGet, http.MethodPut:
		writeResponseBody(w, http.StatusOK, data)
	case http.MethodDelete:
		w.WriteHeader(http.StatusNoContent)
	}
}

func writeResponseBody(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(statusCode)

	_, err = w.Write(res)
	if err != nil {
		log.Printf("error in writing response %v", err)
	}
}
