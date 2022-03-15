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

	// Healthcheck used by uptime-monitoring services, for example Uptime-Kuma
	r.HandleFunc("/health", healthCheck).Methods("GET")

	r.HandleFunc("/debug", AuthRequired(debugInfo)).Methods("GET")
	r.HandleFunc("/switches", AuthRequired(getSwitches)).Methods("GET")
	r.HandleFunc("/token", AuthRequired(updateToken)).Methods("PUT")

	// Handles power request via the connected 433mhz sender
	r.HandleFunc("/power", AuthRequired(setPower)).Methods("POST")

	// Handles bare codes and sends them via the connected 433mhz sender
	r.HandleFunc("/code", AuthRequired(sendCode)).Methods("POST")

	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	r.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
	log.Debug("Initialized Router")
	return r
}
