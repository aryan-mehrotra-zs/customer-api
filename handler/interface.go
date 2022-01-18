package handler

import "net/http"

type handler interface {
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	DeleteByID(w http.ResponseWriter, r *http.Request)
	UpdateByID(w http.ResponseWriter, r *http.Request)
}
