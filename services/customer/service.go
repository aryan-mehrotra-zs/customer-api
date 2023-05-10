package customer

import (
	"github.com/amehrotra/customer-api/errors"
	"github.com/amehrotra/customer-api/models"
	"github.com/amehrotra/customer-api/services"
	"github.com/amehrotra/customer-api/stores"
)

type service struct {
	store stores.Store
}

func New(store stores.Store) services.Customer {
	return service{store: store}
}

func (s service) Create(c models.Customer) (models.Customer, error) {
	return s.store.Create(c)
}

func (s service) Get(id int) (models.Customer, error) {
	return s.store.Get(id)
}

func (s service) Update(c models.Customer) error {
	if _, err := s.Get(c.ID); err != nil {
		return err
	}

	if err := s.store.Update(c); err != nil {
		return err
	}

	_, err := s.store.Get(c.ID)

	return err
}

func (s service) Delete(id int) error {
	if _, err := s.Get(id); err != nil {
		return errors.EntityNotFound{}
	}

	return s.store.Delete(id)
}
