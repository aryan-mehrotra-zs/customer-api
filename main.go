package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/amehrotra/customer-api/drivers"
	handler "github.com/amehrotra/customer-api/handlers/customer"
	service "github.com/amehrotra/customer-api/services/customer"
	store "github.com/amehrotra/customer-api/stores/customer"
)

func main() {
	db, err := drivers.ConnectToSQL()
	if err != nil {
		return
	}

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	stores := store.New(db)
	services := service.New(stores)
	h := handler.New(services)

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
