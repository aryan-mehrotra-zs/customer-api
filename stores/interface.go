package stores

import "github.com/amehrotra/customer-api/models"

type Store interface {
	Get(id int) (models.Customer, error)
	Create(c models.Customer) (models.Customer, error)
	Delete(id int) error
	Update(c models.Customer) error
}
