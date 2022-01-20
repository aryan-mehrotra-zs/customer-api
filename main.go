package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/amehrotra/customer-api/handler"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/customers/{id}", handler.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/customers", handler.Create).Methods(http.MethodPost)
	r.HandleFunc("/customers/{id}", handler.DeleteByID).Methods(http.MethodDelete)
	r.HandleFunc("/customers/{id}", handler.UpdateByID).Methods(http.MethodPut)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(srv.ListenAndServe())
}
