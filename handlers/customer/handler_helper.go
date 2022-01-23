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
)

func getData(w http.ResponseWriter, r *http.Request) (models.Customer, error) {
	var customer models.Customer

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return models.Customer{}, errors.BindError{}
	}

	err = json.Unmarshal(body, &customer)
	if err != nil {
		return models.Customer{}, errors.InvalidParam{}
	}

	return customer, nil
}

func setStatusCode(w http.ResponseWriter, r *http.Request, err error, data interface{}) {
	switch err.(type) {
	case errors.EntityAlreadyExists:
		w.WriteHeader(http.StatusOK)
	case errors.InvalidParam, errors.MissingParam:
		w.WriteHeader(http.StatusBadRequest)
	case errors.EntityNotFound:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		writeSuccessResponse(w, r.Method, data)
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
		writeResponseBody(w, http.StatusNoContent, data)
	}
}

func writeResponseBody(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(data)
	if err != nil {
		setStatusCode(w, nil, err, nil)

		return
	}

	w.WriteHeader(status)

	_, err = w.Write(resp)
	if err != nil {
		setStatusCode(w, nil, err, nil)
	}

}

func getId(w http.ResponseWriter, r *http.Request) (int, error) {
	param := mux.Vars(r)
	idParam := param["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return -1, errors.InvalidParam{}
	}

	return id, nil
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
		log.Println("error in writing response")
	}
}
