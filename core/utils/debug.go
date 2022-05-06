package utils

import (
	"runtime"

	"github.com/smarthome-go/node/core/config"
)

type Debug struct {
	Version              string          `json:"version"`
	Hardware             config.Hardware `json:"hardware"`
	RunningOnRaspberryPi bool            `json:"runningOnRaspberryPi"`
	CpuCount             uint8           `json:"cpuCount"`
}

func GetDebugInfo() Debug {
	configuration := config.GetConfig()
	return Debug{
		Version:              config.Version,
		Hardware:             configuration.Hardware,
		RunningOnRaspberryPi: runtime.GOARCH == "arm",
		CpuCount:             uint8(runtime.NumCPU()),
	}
}
