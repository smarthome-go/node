package api

import (
	"encoding/json"
	"net/http"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: "not found",
		Error:   "no route matched",
	})
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: "method not allowed",
		Error:   "the requested method is not allowed for this endpoint",
	})
}
