package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/amehrotra/customer-api/driver"
	"github.com/amehrotra/customer-api/model"
)

func GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := driver.ConnectToSQL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer db.Close()

	param := mux.Vars(r)
	id := param["id"]

	var c model.Customer

	err = db.QueryRow("SELECT * FROM customers WHERE ID = ?", id).
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
			log.Println("error in writing resp")
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := driver.ConnectToSQL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

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

	res, err := db.Exec("INSERT INTO customers (id,name,address,phone_no) VALUES (?,?,?,?)", c.ID, c.Name, c.Address, c.PhoneNo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(fmt.Sprintf("Customer with id : %v added to database", id)))
	if err != nil {
		log.Println("Error in writing response")

		return
	}
}

func DeleteByID(w http.ResponseWriter, r *http.Request) {
	db, err := driver.ConnectToSQL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer db.Close()

	param := mux.Vars(r)
	id := param["id"]

	_, err = db.Exec("DELETE FROM customers WHERE id = ?;", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func createPutQuery(id string, c model.Customer) (string, []interface{}) {
	q := make([]string, 0, 3)
	args := make([]interface{}, 0, 4)

	if c.Name != "" {
		q = append(q, " name=?")
		args = append(args, c.Name)
	}

	if c.Address != "" {
		q = append(q, " address=?")
		args = append(args, c.Address)
	}

	if c.PhoneNo != 0 {
		q = append(q, " phone_no=?")
		args = append(args, strconv.Itoa(c.PhoneNo))
	}

	if q == nil {
		return "", args
	}

	args = append(args, id)
	query := "UPDATE customers SET" + strings.Join(q, ",") + " WHERE id = ?;"
	return query, args
}

func UpdateByID(w http.ResponseWriter, r *http.Request) {
	db, err := driver.ConnectToSQL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer db.Close()

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

	param := mux.Vars(r)
	id := param["id"]

	query, args := createPutQuery(id, c)
	if query == "" {
		return
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = db.QueryRow("SELECT * FROM customers WHERE ID = ?", id).
		Scan(&c.ID, &c.Name, &c.Address, &c.PhoneNo)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	res, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = w.Write(res)
	if err != nil {
		log.Println("error in writing resp")
	}
}
