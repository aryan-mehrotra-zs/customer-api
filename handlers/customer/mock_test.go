package customer

import (
	"net/http"

	"github.com/amehrotra/customer-api/errors"
	"github.com/amehrotra/customer-api/models"
)

type mockService struct{}

func (m mockService) Get(id int) (models.Customer, error) {
	switch id {
	case 1:
		return models.Customer{ID: 1, Name: "Aryan", Address: "Patna", PhoneNo: 1}, nil
	case 0:
		return models.Customer{}, errors.InvalidParam{}
	case 3:
		return models.Customer{}, errors.EntityNotFound{}
	case 4:
		return models.Customer{}, errors.Error("some error")
	}

	return models.Customer{}, nil
}

func (m mockService) Create(c models.Customer) (models.Customer, error) {
	switch c.Name {
	case "Umang":
		return models.Customer{}, errors.EntityAlreadyExists{}
	case "Aryan":
		c.ID = 1
		return c, nil
	case "Ruchit":
		return models.Customer{}, errors.MissingParam{}
	case "Aakanksha":
		return models.Customer{}, errors.Error("some error")
	default:
		return models.Customer{}, nil
	}
}

func (m mockService) Update(c models.Customer) (models.Customer, error) {
	switch c.Name {
	case "aakanksha":
		c.ID = 1
		return c, nil
	case "Aryan":
		return models.Customer{}, errors.EntityNotFound{}
	default:
		return models.Customer{}, errors.Error("some error")
	}
}

func (m mockService) Delete(id int) error {
	switch id {
	case 1:
		return nil
	case 10:
		return errors.EntityNotFound{}
	default:
		return errors.Error("some error")
	}
}

type mockReader struct{}

func (m mockReader) Read(p []byte) (n int, err error) {
	return 0, errors.InvalidParam{}
}

type mockResponseWriter struct {
}

func (m mockResponseWriter) Header() http.Header {
	header := make(map[string][]string)

	return header
}

func (m mockResponseWriter) Write([]byte) (int, error) {
	return 0, errors.Error("some error")
}

func (m mockResponseWriter) WriteHeader(statusCode int) {

}
