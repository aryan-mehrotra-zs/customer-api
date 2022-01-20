package services

import "github.com/amehrotra/customer-api/models"

type Service interface {
	Get(id int) (models.Customer, error)
	Create(c models.Customer) (models.Customer, error)
	Update(c models.Customer) error
	Delete(id int) error
}
