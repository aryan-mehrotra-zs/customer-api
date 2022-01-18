package customer

import (
	"github.com/amehrotra/customer-api/model"
	"github.com/amehrotra/customer-api/stores"
)

type service struct {
	store stores.Store
}

func New(store stores.Store) service {
	return service{store: store}
}

func (s service) Get(id int) model.Customer {
	c, _ := s.store.Get(id)

	// todo : how to handle error?
	return c
}

func (s service) Create(c model.Customer) int64 {
	resp, _ := s.store.Create(c)

	//todo : call store get and check if correct values are inserted
	return resp
}

func (s service) Delete(id int) {
	s.store.Delete(id)
}

func (s service) Update(c model.Customer) {
	s.store.Update(c)
}
