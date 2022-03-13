package utils

import "github.com/MikMuellerDev/smarthome-hw/core/config"

type Debug struct {
	Version              string `json:"version"`
	Production           bool   `json:"production"`
	RunningOnRaspberryPi bool   `json:"runningOnRaspberryPi"`
}

func GetDebugInfo() Debug {
	configuration := config.GetConfig()
	return Debug{
		Version:              configuration.Version,
		Production:           configuration.Production,
		RunningOnRaspberryPi: configuration.GOARCH == "arm",
	}
}
