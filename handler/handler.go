package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/amehrotra/customer-api/driver"
	"github.com/amehrotra/customer-api/model"
)

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := driver.ConnectToSQL()
	defer db.Close()

	param := mux.Vars(r)
	id := param["id"]

	var c model.Customer

	err := db.QueryRow("SELECT * FROM customers WHERE ID = ?", id).
		Scan(&c.ID, &c.Name, &c.Address, &c.PhoneNo)

	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		res, err := json.Marshal(c)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)

		_, err = w.Write(res)
		if err != nil {
			log.Println("error in writing response")
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := driver.ConnectToSQL()
	defer db.Close()

	var c model.Customer

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	result, err := db.Exec("INSERT INTO customers (id,name,address,phone_no) VALUES (?,?,?,?)", c.ID, c.Name, c.Address, c.PhoneNo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	log.Println(result.LastInsertId())

	w.WriteHeader(http.StatusCreated)
}
