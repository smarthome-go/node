package firmware

import (
	"fmt"

	"github.com/smarthome-go/node/core/config"
	"github.com/smarthome-go/node/core/log"
	"github.com/smarthome-go/rpirf"
)

var sender rpirf.RFDevice

func Init(config config.Hardware) error {
	var output rpirf.HardwareOutput

	if config.UseGpioCdev {
		outputTemp, err := rpirf.NewCharacterDev(
			*config.GpioCdevChip,
			config.RFDevicePin,
		)

		if err != nil {
			log.Error("Failed to initialize CDEV RFDevice: ", err.Error())
			return err
		}

		output = outputTemp
	} else {
		outputTemp, err := rpirf.NewRaspberryPi(
			uint(config.RFDevicePin),
		)

		if err != nil {
			log.Error("Failed to initialize RaspberryPi RFDevice: ", err.Error())
			return err
		}

		output = outputTemp
	}

	// Initialize the sender
	sender = rpirf.NewRF(
		output,
		config.RFDeviceProtocol,
		config.RFDeviceRepeat,
		config.RFDevicePulselength,
		config.RFDeviceLength,
	)

	log.Trace(fmt.Sprintf("RFDevice (433mhz) on pin %d initialized. RFDevice repeat: %d", config.RFDevicePin, config.RFDeviceRepeat))
	return nil
}

func Cleanup() error {
	return sender.Cleanup()
}

func sendCode(
	code int,
	config config.Hardware,
) error {
	// Send the code
	return sender.Send(code)
}
