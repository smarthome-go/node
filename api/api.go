package api

import (
	"encoding/json"
	"net/http"

	"github.com/smarthome-go/node/core/config"
	"github.com/smarthome-go/node/core/firmware"
	"github.com/smarthome-go/node/core/log"
	"github.com/smarthome-go/node/core/utils"
)

type UpdateTokenRequest struct {
	Token string `json:"token"`
}

type PowerRequest struct {
	Switch  string `json:"switch"`
	PowerOn bool   `json:"power"`
}

type CodeRequest struct {
	Code int `json:"code"`
}

// Used when requesting a list of available switches
type SwitchListing struct {
	SwitchesRF   []config.SwitchRF   `json:"switchesRF"`
	SwitchesGPIO []config.SwitchGPIO `json:"switchesGPIO"`
}

// Returns general-purpose debugging information
func debugInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(utils.GetDebugInfo()); err != nil {
		log.Error("Encoding json failed: ", err.Error())
	}
}

// Returns a service unavailable if the hardware is disabled
func healthCheck(w http.ResponseWriter, r *http.Request) {
	if config.GetConfig().Hardware.HardwareEnabled {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

// Can be used to update the token and generate a new hash
func updateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request UpdateTokenRequest
	if err := decoder.Decode(&request); err != nil || request.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(Response{Success: false, Message: "bad request", Error: "invalid request body"}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	}
	hash, err := config.GenerateTokenHash(request.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "failed to update token",
			Error:   "could not generate token",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	}
	config.SetHash(hash)
	if err := config.WriteConfig(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "failed to update token",
			Error:   "could not write generated token to config file",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	}
	if err := json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "successfully updated token",
	}); err != nil {
		log.Error("Encoding json failed: ", err.Error())
	}
}

// Returns all switches, also their internal 433mhz codes and GPIO pin numberrings
func getSwitches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(SwitchListing{
		SwitchesRF:   config.GetConfig().SwitchesRF,
		SwitchesGPIO: config.GetConfig().SwitchesGPIO,
	}); err != nil {
		log.Error("Encoding json failed: ", err.Error())
	}
}

// Main function used to communicate with the hardware
func setPower(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Decode the request body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request PowerRequest
	if err := decoder.Decode(&request); err != nil || request.Switch == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "bad request",
			Error:   "invalid request body",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	}
	err := firmware.SetPower(
		request.Switch,
		request.PowerOn,
	)

	if err != nil {
		log.Error(err.Error())
	}

	switch err {
	case firmware.ErrBlocked:
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "currently blocked",
			Error:   "the sender's hardware is currently busy",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	case firmware.ErrDisabled:
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "currently disabled / maintenance mode",
			Error:   "the hardware is currently disabled",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	case firmware.ErrNoSwitch:
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "bad request",
			Error:   "the desired switch could not be matched to a code, is it valid?",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	case nil:
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "power request successful",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	default:
		log.Error("Failed to match error return value for power request: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "internal server error",
			Error:   "could not match return value of hardware power request",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	}
}

// Accept a bare code and sends it (used for testing)
func sendCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Decode the request body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request CodeRequest
	if err := decoder.Decode(&request); err != nil || request.Code == 0 {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "bad request",
			Error:   "invalid request body",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	}
	err := firmware.SendCode(request.Code)
	switch err {
	case firmware.ErrBlocked:
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "currently blocked",
			Error:   "the sender's hardware is currently busy",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	case firmware.ErrDisabled:
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "currently disabled / maintenance mode",
			Error:   "the hardware is currently disabled",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	case nil:
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "power request successful",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	default:
		log.Error("Failed to match error return value for power request: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "internal server error",
			Error:   "could not match return value of hardware power request",
		}); err != nil {
			log.Error("Encoding json failed: ", err.Error())
		}
		return
	}
}
