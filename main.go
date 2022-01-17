package main

import (
	"database/sql"
	"fmt"
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
var err error

func main() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "organisation",
	}
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Println(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Println(pingErr)
	}
	fmt.Println("Connected!")

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
	var c customer

	q := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT * FROM customers WHERE ID = ?", q)

	err := row.Scan(&c.ID, &c.Name, &c.Address, &c.PhoneNo)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(c)

}
