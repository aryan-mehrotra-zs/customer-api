package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/amehrotra/customer-api/driver"
	"github.com/amehrotra/customer-api/model"
)

func GetByID(w http.ResponseWriter, r *http.Request) {
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

	id, _ := result.LastInsertId()
	log.Printf("customer added with id %v\n", id)

	w.WriteHeader(http.StatusCreated)
}

func DeleteByID(w http.ResponseWriter, r *http.Request) {
	db := driver.ConnectToSQL()
	defer db.Close()

	param := mux.Vars(r)
	id := param["id"]

	_, err := db.Exec("DELETE FROM customers WHERE id = ?;", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func Put(w http.ResponseWriter, r *http.Request) {
	db := driver.ConnectToSQL()
	defer db.Close()

	param := mux.Vars(r)
	id := param["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var c model.Customer
	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	query := "UPDATE customers SET"

	var q []string

	if c.Name != "" {
		q = append(q, " name = \""+c.Name+"\"")
	}
	if c.Address != "" {
		q = append(q, " address = \""+c.Address+"\"")
	}
	if c.PhoneNo != 0 {
		q = append(q, " phone_no = "+string(c.PhoneNo))
	}

	if q == nil {

		return
	}

	query += strings.Join(q, ",")
	query += " WHERE id=" + string(id) + ";"
	db.Exec(query)
}
