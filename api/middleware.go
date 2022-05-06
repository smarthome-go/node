package api

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/smarthome-go/node/core/config"
	"github.com/smarthome-go/node/core/log"
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
			if err := json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "unauthorized",
				Error:   "no token provided",
			}); err != nil {
				log.Error("Encoding json failed: ", err.Error())
			}
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(config.GetConfig().TokenHash), []byte(token)); err == nil {
			handler.ServeHTTP(w, r)
			return
		}
		log.Debug("invalid token provided, rejecting request")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "unauthorized",
			Error:   "invalid credentials",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
	}
}
