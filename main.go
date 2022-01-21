package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/amehrotra/customer-api/drivers"
	handler "github.com/amehrotra/customer-api/handlers/customer"
	store "github.com/amehrotra/customer-api/stores/customer"
)

func main() {
	db := drivers.ConnectToSQL()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	stores := store.New(db)
	h := handler.New(stores)

	r := mux.NewRouter()

	r.HandleFunc("/customers/{id}", h.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/customers", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/customers/{id}", h.DeleteByID).Methods(http.MethodDelete)
	r.HandleFunc("/customers/{id}", h.UpdateByID).Methods(http.MethodPut)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatalln(srv.ListenAndServe())
}
