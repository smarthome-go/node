package firmware

import (
	"fmt"

	"github.com/smarthome-go/node/core/config"
	"github.com/smarthome-go/node/core/log"
	"github.com/smarthome-go/rpirf"
)

func sendCode(
	code int,
	config config.Hardware,
) error {
	// Initialize the sender
	sender, err := rpirf.NewRF(
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
	log.Trace(fmt.Sprintf("RFDevice (433mhz) on pin %d initialized. RFDevice repeat: %d", config.RFDevicePin, config.RFDeviceRepeat))
	// Send the code
	if err := sender.Send(code); err != nil {
		return err
	}
	// Free the sender's GPIO device
	return sender.Cleanup()
}

func TestSender(config config.Hardware) error {
	// Initialize the sender
	sender, err := rpirf.NewRF(
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
	log.Info(fmt.Sprintf("RFDevice (433mhz) on pin %d initialized. RFDevice repeat: %d", config.RFDevicePin, config.RFDeviceRepeat))
	// Free the sender's GPIO device
	return sender.Cleanup()
}
