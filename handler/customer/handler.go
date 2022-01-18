package customer

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/amehrotra/customer-api/model"
	"github.com/amehrotra/customer-api/services"
)

type handler struct {
	service services.Service
}

func New(service services.Service) handler {
	return handler{service: service}
}

func (h handler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param := mux.Vars(r)
	idParam := param["id"]
	id, err := strconv.Atoi(idParam)

	data, err := h.service.Get(id)

	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
	case nil:
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

	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c model.Customer

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	c, err = h.service.Create(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	idParam := param["id"]
	id, err := strconv.Atoi(idParam)

	err = h.service.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h handler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var c model.Customer

	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	param := mux.Vars(r)
	idParam := param["id"]
	c.ID, err = strconv.Atoi(idParam)

	err = h.service.Update(c)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}
