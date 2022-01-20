package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetByID(t *testing.T) {
	cases := []struct {
		desc       string
		id         string // body
		resp       []byte
		statusCode int
	}{
		{"id exists in db", "4", []byte(`{"id":4,"name":"Umang","address":"India","phone_no":4}`), http.StatusOK},
	}

	for i, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, "http://dummy", http.NoBody)
		r := mux.SetURLVars(req, map[string]string{"id": tc.id})
		w := httptest.NewRecorder()

		GetByID(w, r)

		resp := w.Result()
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Errorf("cannot read resp: %v", err)
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, resp.StatusCode, tc.statusCode)
		}

		if !reflect.DeepEqual(body, tc.resp) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, string(body), string(tc.resp))
		}
	}
}

func TestCreate(t *testing.T) {
	cases := []struct {
		desc       string
		body       []byte
		statusCode int
	}{
		{"already exists parameter", []byte(`{"id":4,"name":"Umang","address":"India","phone_no":4}`), http.StatusInternalServerError},
		{"create new customer", []byte(`{"name":"Umang","address":"India","phone_no":4}`), http.StatusCreated},
	}

	for i, tc := range cases {
		r := httptest.NewRequest(http.MethodPost, "http://dummy", bytes.NewReader(tc.body))
		w := httptest.NewRecorder()

		Create(w, r)

		resp := w.Result()

		if resp.StatusCode != tc.statusCode {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, resp.StatusCode, tc.statusCode)
		}
	}
}

func TestDeleteByID(t *testing.T) {
	cases := []struct {
		desc       string
		id         string
		statusCode int
	}{
		{"deleted successful", "3", http.StatusNoContent},
	}

	for i, tc := range cases {
		req := httptest.NewRequest(http.MethodDelete, "http://dummy", http.NoBody)
		r := mux.SetURLVars(req, map[string]string{"id": tc.id})
		w := httptest.NewRecorder()

		DeleteByID(w, r)

		resp := w.Result()

		if resp.StatusCode != tc.statusCode {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, resp.StatusCode, tc.statusCode)
		}

	}
}

func TestUpdateByID(t *testing.T) {
	cases := []struct {
		desc       string
		id         string
		body       []byte
		customer   []byte
		statusCode int
	}{
		{"Successfully updated", "8", []byte(`{"phone_no":6}`), []byte(`{"id":8,"name":"Umang","address":"India","phone_no":6}`), http.StatusCreated},
	}

	for i, tc := range cases {
		req := httptest.NewRequest(http.MethodPut, "http://dummy", bytes.NewReader(tc.body))
		r := mux.SetURLVars(req, map[string]string{"id": tc.id})
		w := httptest.NewRecorder()

		UpdateByID(w, r)

		resp := w.Result()

		if resp.StatusCode != tc.statusCode {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, resp.StatusCode, tc.statusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("cannot read resp: %v", err)
		}

		if !reflect.DeepEqual(body, tc.customer) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, string(body), string(tc.customer))
		}
	}
}
