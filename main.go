package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type customer struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Address string `json:"Address"`
	PhoneNo int    `json:"PhoneNo"`
}

var db *sql.DB

func ConnectToSQL() *sql.DB {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "organisation",
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Println(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
	}
	log.Println("Connected!")

	return db
}

func main() {
	db = ConnectToSQL()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/customers/{id}", Get).Methods(http.MethodGet)
	r.HandleFunc("/customers", Post).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(srv.ListenAndServe())
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c customer

	parameter := mux.Vars(r)
	id := parameter["id"]

	row := db.QueryRow("SELECT * FROM customers WHERE ID = ?", id)
	err := row.Scan(&c.ID, &c.Name, &c.Address, &c.PhoneNo)

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
		w.Write(res)

	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var c customer

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	result, err := db.Exec("INSERT INTO customers (ID,Name,Address,PhoneNo) VALUES (?,?,?,?)", c.ID, c.Name, c.Address, c.PhoneNo)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	log.Println(result.LastInsertId())

}

/*
INSERT INTO customers (ID,Name,Address,PhoneNo) VALUES (1,"Aryan","Patna",9852902205);


*/
