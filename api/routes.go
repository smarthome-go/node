package api

import (
	"net/http"

	"github.com/MikMuellerDev/smarthome-hw/core/log"
	"github.com/gorilla/mux"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/debug", debugInfo).Methods("GET")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	r.HandleFunc("/token", AuthRequired(updateToken)).Methods("PUT")

	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	r.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
	log.Debug("Initialized Router.")
	return r
}
