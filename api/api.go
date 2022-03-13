package api

import (
	"encoding/json"
	"net/http"

	"github.com/MikMuellerDev/smarthome-hw/core/utils"
)

func debugInfo(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(utils.GetDebugInfo())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func updateToken(w http.ResponseWriter, r *http.Request) {

}
