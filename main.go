package main

import (
	"database/sql"
	"encoding/json"
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

	r.HandleFunc("/customers", Get).Methods(http.MethodGet)

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

	q := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT * FROM customers WHERE ID = ?", q)

	err := row.Scan(&c.ID, &c.Name, &c.Address, &c.PhoneNo)
	if err != nil {
		log.Println(err.Error())
	}

	res, err := json.Marshal(c)
	w.Write(res)

}
