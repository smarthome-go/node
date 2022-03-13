package main

import (
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
	config.SetVersion("v0.0.1")
}
