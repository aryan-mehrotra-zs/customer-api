package customer

import (
	"bytes"
	"github.com/amehrotra/customer-api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockService struct {
}

func (m mockService) Get(id int) (models.Customer, error) {

}

func (m mockService) Create(c models.Customer) (models.Customer, error) {
	switch c {
	case models.Customer{
		ID:      4,
		Name:    "Umang",
		Address: "India",
		PhoneNo: 4,
	}:
		return models.Customer{}, errors.EntityAlreadyExists{}
	}
}

func (m mockService) Update(c models.Customer) error {

}
func (m mockService) Delete(id int) error {

}

func TestHandler_Create(t *testing.T) {
	h := New(mockService{})

	cases := []struct {
		desc       string
		body       []byte
		statusCode int
	}{
		{"already exists parameter", []byte(`{"id":4,"name":"Umang","address":"India","phone_no":4}`), http.StatusInternalServerError},
		{"create new resp", []byte(`{"name":"Umang","address":"India","phone_no":4}`), http.StatusCreated},
	}

	for i, tc := range cases {
		r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))
		w := httptest.NewRecorder()

		h.Create(w, r)

		resp := w.Result()

		if resp.StatusCode != tc.statusCode {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, resp.StatusCode, tc.statusCode)
		}
	}
}
