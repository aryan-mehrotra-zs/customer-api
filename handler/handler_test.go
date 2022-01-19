package handler

import (
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
		id         string // input
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
