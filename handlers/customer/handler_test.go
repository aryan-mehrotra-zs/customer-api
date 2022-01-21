package customer

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/amehrotra/customer-api/errors"
	"github.com/amehrotra/customer-api/models"
)

type mockService struct {
}

func (m mockService) Get(id int) (models.Customer, error) {
	switch strconv.Itoa(id) {
	case "1":
		return models.Customer{ID: 1, Name: "Aryan", Address: "Patna", PhoneNo: 1}, nil
	case "0":
		return models.Customer{}, errors.InvalidParam{}
	case "3":
		return models.Customer{}, errors.EntityNotFound{}
	case "4":
		return models.Customer{}, errors.DB{}
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
		return models.Customer{}, errors.DB{}
	default:
		return models.Customer{}, nil
	}
}

func (m mockService) Update(c models.Customer) error {
	return nil
}
func (m mockService) Delete(id int) error {
	return nil
}

type mockReader struct{}

func (m mockReader) Read(p []byte) (n int, err error) {
	return 0, errors.InvalidParam{}
}

func TestHandler_Create(t *testing.T) {
	h := New(mockService{})

	cases := []struct {
		desc       string
		body       io.Reader
		resp       []byte
		statusCode int
	}{
		{"entity already exists", bytes.NewReader([]byte(`{"id":4,"name":"Umang"}`)), nil, http.StatusOK},
		{"create new entity", bytes.NewReader([]byte(`{"name":"Aryan"}`)), nil, http.StatusCreated},
		{"missing or invalid parameter", bytes.NewReader([]byte(`{"name":"Ruchit"}`)), nil, http.StatusBadRequest},
		{"internal server error", bytes.NewReader([]byte(`{"name":"Aakanksha"}`)), nil, http.StatusInternalServerError},
		{"unmarshal error", bytes.NewReader([]byte(`invalid body`)), nil, http.StatusBadRequest},
		{"bind error", mockReader{}, nil, http.StatusInternalServerError},
	}

	for i, tc := range cases {
		r := httptest.NewRequest(http.MethodPost, "http://dummy", tc.body)
		w := httptest.NewRecorder()

		h.Create(w, r)

		resp := w.Result()

		err := resp.Body.Close()
		if err != nil {
			t.Errorf("error in writing response")
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}
	}
}

func TestHandler_GetByID(t *testing.T) {
	h := New(mockService{})

	cases := []struct {
		desc       string
		id         string
		resp       []byte
		statusCode int
	}{
		{"entity exists", "1", []byte(`{"id":1,"name":"Aryan","address":"Patna","phone_no":1}`), http.StatusOK},
		{"invalid id", "0", []byte(``), http.StatusBadRequest},
		{"entity not found", "3", []byte(``), http.StatusNotFound},
		{"database connectivity error", "4", []byte(``), http.StatusInternalServerError},
	}
	for i, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, "http://dummy", http.NoBody)
		r := mux.SetURLVars(req, map[string]string{"id": tc.id})
		w := httptest.NewRecorder()

		h.GetByID(w, r)

		resp := w.Result()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("cannot read resp: %v", err)
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if !reflect.DeepEqual(body, tc.resp) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), string(tc.resp))
		}
	}
}

func TestHandler_writeResponse(t *testing.T) {

	c := make(chan int)
	w := httptest.NewRecorder()
	writeResponse(w, c)
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("\n[TEST] Failed. Desc : Marshal Error\nGot %v\nExpected %v", resp.StatusCode, http.StatusInternalServerError)
	}

}
