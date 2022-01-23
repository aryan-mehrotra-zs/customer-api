package customer

import (
	"net/http"

	"github.com/amehrotra/customer-api/services"
)

type handler struct {
	service services.Customer
}

func New(service services.Customer) handler {
	return handler{service: service}
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	customer, err := getData(r)
	if err != nil {
		setStatusCode(w, r, err, customer)

		return
	}

	customer, err = h.service.Create(customer)

	setStatusCode(w, r, err, customer)
}

func (h handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		setStatusCode(w, r, err, nil)

		return
	}

	customer, err := h.service.Get(id)

	setStatusCode(w, r, err, customer)
}

func (h handler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	customer, err := getData(r)
	if err != nil {
		setStatusCode(w, r, err, nil)

		return
	}

	customer.ID, err = getId(r)
	if err != nil {
		setStatusCode(w, r, err, nil)

		return
	}

	customer, err = h.service.Update(customer)

	setStatusCode(w, r, err, customer)
}

func (h handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		setStatusCode(w, r, err, nil)

		return
	}

	err = h.service.Delete(id)

	setStatusCode(w, r, err, nil)
}
