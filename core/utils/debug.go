package utils

import (
	"runtime"

	"github.com/MikMuellerDev/smarthome-hw/core/config"
)

type Debug struct {
	Version              string `json:"version"`
	Production           bool   `json:"production"`
	RunningOnRaspberryPi bool   `json:"runningOnRaspberryPi"`
}

func GetDebugInfo() Debug {
	configuration := config.GetConfig()
	return Debug{
		Version:              config.Version,
		Production:           configuration.Production,
		RunningOnRaspberryPi: runtime.GOARCH == "arm",
	}
}
