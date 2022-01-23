package customer

import (
	"testing"

	"github.com/amehrotra/customer-api/errors"
	"github.com/amehrotra/customer-api/models"
)

type mockStore struct{}

func (m mockStore) Create(c models.Customer) (models.Customer, error) {
	switch c.Name {
	case "Aryan":
		return models.Customer{1, "Aryan", "Patna", 1}, nil
	case "Umang":
		return models.Customer{}, errors.EntityNotFound{}
	}

	return models.Customer{}, nil
}

func (m mockStore) Get(id int) (models.Customer, error) {
	switch id {
	case 1:
		return models.Customer{1, "Aryan", "Patna", 1}, nil
	case 2:
		return models.Customer{}, errors.EntityNotFound{}
	}
	return models.Customer{}, nil
}

func (m mockStore) Update(c models.Customer) error {
	switch c.ID {
	case 1:
		return errors.EntityNotFound{}
	case 2:
		return errors.InvalidParam{}
	default:
		return nil
	}
}

func (m mockStore) Delete(id int) error {
	switch id {
	case 2:
		return errors.EntityNotFound{}
	default:
		return nil
	}
}

func TestService_Create(t *testing.T) {
	s := New(mockStore{})

	cases := []struct {
		desc     string
		customer models.Customer
		resp     models.Customer
		err      error
	}{
		{"create successful",
			models.Customer{
				Name:    "Aryan",
				Address: "Patna",
				PhoneNo: 1,
			}, models.Customer{
			ID:      1,
			Name:    "Aryan",
			Address: "Patna",
			PhoneNo: 1,
		}, nil},
		{"create unsuccessful", models.Customer{
			Name:    "Umang",
			Address: "Patna",
			PhoneNo: 1,
		}, models.Customer{}, errors.EntityNotFound{}},
	}

	for i, tc := range cases {
		resp, err := s.Create(tc.customer)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if tc.resp != resp {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_Get(t *testing.T) {
	s := New(mockStore{})

	cases := []struct {
		desc string
		id   int
		resp models.Customer
		err  error
	}{
		{"entity found successful", 1, models.Customer{
			ID:      1,
			Name:    "Aryan",
			Address: "Patna",
			PhoneNo: 1,
		}, nil},
		{"entity not found", 2, models.Customer{}, errors.EntityNotFound{}},
	}

	for i, tc := range cases {
		resp, err := s.Get(tc.id)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}

		if tc.resp != resp {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.resp)
		}
	}
}

func TestService_Update(t *testing.T) {
	s := New(mockStore{})

	cases := []struct {
		desc     string
		customer models.Customer
		err      error
	}{
		{"entity does not exist", models.Customer{
			ID:      1,
			Name:    "Aryan",
			Address: "Patna",
			PhoneNo: 1,
		}, errors.EntityNotFound{}},
		{"successfully updated", models.Customer{
			ID:      3,
			Name:    "Aryan",
			Address: "Patna",
			PhoneNo: 1,
		}, nil},
	}

	for i, tc := range cases {
		err := s.Update(tc.customer)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\nExpected %v", i, tc.desc, err, tc.err)
		}
	}
}

func TestService_Delete(t *testing.T) {
	s := New(mockStore{})

	cases := []struct {
		desc string
		id   int
		err  error
	}{
		{"Successful Delete", 1, nil},
		{"entity does not exist", 2, errors.EntityNotFound{}},
	}

	for i, tc := range cases {
		err := s.Delete(tc.id)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\nExpected %v", i, tc.desc, err, tc.err)
		}
	}
}
