package customer

import (
	"github.com/amehrotra/customer-api/models"
	"github.com/amehrotra/customer-api/services"
	"github.com/amehrotra/customer-api/stores"
)

type service struct {
	store stores.Store
}

// New fixme why it is not being used
func New(store stores.Store) services.Customer {
	return service{store: store}
}

func (s service) Get(id int) (models.Customer, error) {
	// todo : how to handle error?
	return s.store.Get(id)
}

func (s service) Create(c models.Customer) (models.Customer, error) {
	// todo : call store get and check if correct values are inserted
	return s.store.Create(c)
}

func (s service) Update(c models.Customer) (models.Customer, error) {
	if _, err := s.Get(c.ID); err != nil {
		return models.Customer{}, err
	}

	if err := s.store.Update(c); err != nil {
		return models.Customer{}, err
	}

	return s.store.Get(c.ID)
}

func (s service) Delete(id int) error {
	return s.store.Delete(id)
}
