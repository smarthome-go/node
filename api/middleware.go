package api

import (
	"encoding/json"
	"net/http"

	"github.com/MikMuellerDev/smarthome-hw/core/config"
	"github.com/MikMuellerDev/smarthome-hw/core/log"
	"golang.org/x/crypto/bcrypt"
)

/*
Middleware:
The following function will server the purpose of validating request
Validation includes checking if the provided credentials are valid
Due to all routes returning JSON, the response header is already set here
*/

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		query := r.URL.Query()
		token := query.Get("token")
		if token == "" {
			log.Debug("No authentication token provided, rejecting request")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "unauthorized",
				Error:   "no token provided",
			})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(config.GetConfig().TokenHash), []byte(token)); err == nil {
			handler.ServeHTTP(w, r)
			return
		}
		log.Debug("invalid token provided, rejecting request")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "unauthorized",
			Error:   "invalid credentials",
		})
	}
}
