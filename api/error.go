package api

import (
	"encoding/json"
	"net/http"

	"github.com/MikMuellerDev/smarthome-hw/core/log"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: "not found",
		Error:   "no route matched",
	}); err != nil {
		log.Error("Encoding json failed: ", err.Error())
	}
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	if err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: "method not allowed",
		Error:   "the requested method is not allowed for this endpoint",
	}); err != nil {
		log.Error("Encoding json failed: ", err.Error())
	}
}
