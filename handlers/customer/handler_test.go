package customer

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func initializeTest(method string, body io.Reader, pathParams map[string]string) (handler, *http.Request, *httptest.ResponseRecorder) {
	h := New(mockService{})

	req := httptest.NewRequest(method, "http://customers", body)
	r := mux.SetURLVars(req, pathParams)
	w := httptest.NewRecorder()

	return h, r, w
}

func TestHandler_Create(t *testing.T) {
	cases := []struct {
		desc       string
		body       io.Reader
		resp       []byte
		statusCode int
	}{
		{"entity already exists", bytes.NewReader([]byte(`{"id":4,"name":"Umang"}`)), nil, http.StatusOK},
		{"create new entity", bytes.NewReader([]byte(`{"name":"Aryan"}`)), nil, http.StatusCreated},
		{"missing or invalid parameter", bytes.NewReader([]byte(`{"name":"Ruchit"}`)), nil, http.StatusBadRequest},
		{"default case : internal server error", bytes.NewReader([]byte(`{"name":"Aakanksha"}`)), nil, http.StatusInternalServerError},
		{"unmarshal error", bytes.NewReader([]byte(`invalid body`)), nil, http.StatusBadRequest},
		{"unable to read body", mockReader{}, nil, http.StatusInternalServerError},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodPost, tc.body, nil)

		h.Create(w, r)

		resp := w.Result()

		body, err := getResponseBody(resp)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if !reflect.DeepEqual(body, tc.resp) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), string(tc.resp))
		}
	}
}

func TestHandler_GetByID(t *testing.T) {
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
		h, r, w := initializeTest(http.MethodGet, http.NoBody, map[string]string{"id": tc.id})

		h.GetByID(w, r)

		resp := w.Result()

		body, err := getResponseBody(resp)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if !reflect.DeepEqual(body, tc.resp) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), string(tc.resp))
		}
	}
}

func TestHandler_UpdateByID(t *testing.T) {
	cases := []struct {
		desc       string
		id         string
		body       io.Reader
		statusCode int
		resp       []byte
	}{
		{"entity updated successfully", "1", bytes.NewReader([]byte(`{"name":"aakanksha","address":"Patna","phone_no":1}`)), http.StatusOK, []byte(`{"id":1,"name":"aakanksha","address":"Patna","phone_no":1}`)},
		{"entity not found", "10", bytes.NewReader([]byte(`{"name":"Aryan"}`)), http.StatusNotFound, []byte("")},
		{"server error", "99", bytes.NewReader([]byte(`{"name":"Umang"}`)), http.StatusInternalServerError, []byte("")},
		{"invalid id", "abc", bytes.NewReader([]byte(`{"name":"Umang"}`)), http.StatusBadRequest, []byte("")},
		{"unmarshal error", "10", bytes.NewReader([]byte(`invalid body`)), http.StatusBadRequest, []byte("")},
		{"body read error", "10", mockReader{}, http.StatusInternalServerError, []byte("")},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodPut, tc.body, map[string]string{"id": tc.id})

		h.UpdateByID(w, r)

		resp := w.Result()

		body, err := getResponseBody(resp)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}

		if !reflect.DeepEqual(body, tc.resp) {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, string(body), string(tc.resp))
		}
	}
}

func TestHandler_DeleteByID(t *testing.T) {
	cases := []struct {
		desc       string
		id         string
		statusCode int
	}{
		{"delete successful", "1", http.StatusNoContent},
		{"entity not found", "10", http.StatusNotFound},
		{"server error", "11", http.StatusInternalServerError},
		{"invalid id", "abc", http.StatusBadRequest},
	}

	for i, tc := range cases {
		h, r, w := initializeTest(http.MethodDelete, http.NoBody, map[string]string{"id": tc.id})

		h.DeleteByID(w, r)

		resp := w.Result()

		if resp.StatusCode != tc.statusCode {
			t.Errorf("\n[TEST %d] Failed. Desc : %v\nGot %v\nExpected %v", i, tc.desc, resp.StatusCode, tc.statusCode)
		}
	}
}

func TestHandler_writeResponseMarshalError(t *testing.T) {
	data := make(chan int)
	w := httptest.NewRecorder()
	expectedStatusCode := http.StatusInternalServerError

	writeResponseBody(w, data)

	resp := w.Result()
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("\n[TEST] Failed. Desc : Marshal Error \nGot %v\nExpected %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestHandler_writeResponseWriteError(t *testing.T) {
	data := []byte(`{"id":1,"name":"Aryan","address":"Patna","phone_no":1}`)
	w := mockResponseWriter{}

	var b bytes.Buffer
	log.SetOutput(&b)

	writeResponseBody(w, data)

	if !strings.Contains(b.String(), "error in writing response") {
		t.Errorf("\n[TEST] Failed. Desc : Write Error \nGot %v\nExpected 'error in writing response' in logs", b.String())
	}
}

func getResponseBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil
}
