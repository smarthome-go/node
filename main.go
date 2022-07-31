package main

import (
	"fmt"
	"net/http"

	"github.com/smarthome-go/node/api"
	"github.com/smarthome-go/node/core/config"
	"github.com/smarthome-go/node/core/firmware"
	"github.com/smarthome-go/node/core/log"
)

func main() {
	if err := log.InitLogger(6); err != nil {
		panic(err.Error())
	}
	if err := config.ReadConfigFile(); err != nil {
		log.Fatal("Failed to read config file: ", err.Error())
	}
	config.Version = "0.3.0"
	log.Debug("Successfully read config file")

	r := api.NewRouter()
	http.Handle("/", r)

	if !config.GetConfig().Hardware.HardwareEnabled {
		log.Warn("Hardware is disabled, this server will not works as intended")
	} else {
		// If the hardware is enabled and the software is run on a raspberry pi (arm), enable the sender
		if err := firmware.Init(); err != nil {
			log.Warn("Deactivating hardware due to previous initialization failure")
			config.SetHardwareEnabled(false)
			if err := config.WriteConfig(); err != nil {
				log.Fatal("Failed to deactivate hardware after initialization failure", err.Error())
			}
		}
		defer func() {
			if err := firmware.Free(); err != nil {
				log.Fatal("Could not deactivate sender: ", err.Error())
			}
		}()
	}
	log.Info(fmt.Sprintf("Smarthome-hw %s is running on http://localhost:%d", config.Version, config.GetConfig().Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.GetConfig().Port), r).Error())
}
