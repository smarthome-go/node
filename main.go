package main

import (
	"fmt"
	"net/http"

	"github.com/MikMuellerDev/smarthome-hw/api"
	"github.com/MikMuellerDev/smarthome-hw/core/config"
	"github.com/MikMuellerDev/smarthome-hw/core/log"
)

func main() {
	if err := log.InitLogger(6); err != nil {
		panic(err.Error())
	}
	if err := config.ReadConfigFile(); err != nil {
		log.Fatal("Failed to read config file: ", err.Error())
	}
	config.Version = "v0.0.1"
	r := api.NewRouter()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.GetConfig().Port), r).Error())
}
