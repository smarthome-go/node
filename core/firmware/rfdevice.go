package firmware

import (
	"fmt"

	"github.com/smarthome-go/node/core/config"
	"github.com/smarthome-go/node/core/log"
	"github.com/smarthome-go/rpirf"
)

var sender rpirf.RFDevice

func Init() error {
	config := config.GetConfig().Hardware
	device, err := rpirf.NewRF(
		config.RFDevicePin,
		config.RFDeviceProtocol,
		config.RFDeviceRepeat,
		config.RFDevicePulselength,
		config.RFDeviceLength,
	)
	if err != nil {
		log.Error("Failed to initialize RFDevice: ", err.Error())
		return err
	}
	sender = device
	log.Info(fmt.Sprintf("RFDevice (433mhz) on pin %d initialized. RFDevice repeat: %d", config.RFDevicePin, config.RFDeviceRepeat))
	return nil
}

func Free() error {
	if err := sender.Cleanup(); err != nil {
		log.Error("Failed to free sender: ", err.Error())
		return err
	}
	return nil
}
