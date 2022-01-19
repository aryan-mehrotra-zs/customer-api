package customer

import (
	"github.com/amehrotra/customer-api/model"
	"github.com/amehrotra/customer-api/stores"
)

type service struct {
	store stores.Store
}

// New fixme why it is not being used
func New(store stores.Store) service {
	return service{store: store}
}

func (s service) Get(id int) (model.Customer, error) {

	// todo : how to handle error?
	return s.store.Get(id)
}

func (s service) Create(c model.Customer) (model.Customer, error) {
	//todo : call store get and check if correct values are inserted
	return s.store.Create(c)
}

func (s service) Update(c model.Customer) error {
	return s.store.Update(c)
}

func (s service) Delete(id int) error {
	return s.store.Delete(id)
}
