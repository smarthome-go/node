package api

import (
	"encoding/json"
	"net/http"

	"github.com/MikMuellerDev/smarthome-hw/core/config"
	"github.com/MikMuellerDev/smarthome-hw/core/utils"
)

type UpdateTokenRequest struct {
	Token string `json:"token"`
}

func debugInfo(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(utils.GetDebugInfo())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func updateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request UpdateTokenRequest
	if err := decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: false, Message: "bad request", Error: "invalid request body"})
		return
	}
	hash, err := config.GenerateTokenHash(request.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "failed to update token",
			Error:   "could not generate token",
		})
		return
	}
	config.SetHash(hash)
	if err := config.WriteConfig(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "failed to update token",
			Error:   "could not write generated token to config file",
		})
		return
	}
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "successfully updated token",
	})
}
