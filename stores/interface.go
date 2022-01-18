package stores

import "github.com/amehrotra/customer-api/model"

type Store interface {
	Get(id int) (model.Customer, error)
	Create(c model.Customer) (model.Customer, error)
	Delete(id int) error
	Update(c model.Customer) error
}
