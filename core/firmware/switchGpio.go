package firmware

import (
	"github.com/smarthome-go/node/core/log"
	"github.com/stianeikeland/go-rpio/v4"
)

func setPin(
	pinNumber uint8,
	powerOn bool,
) error {
	if err := rpio.Open(); err != nil {
		log.Error("Could not set power of GPIO pin: failed to open /dev/mem: ", err.Error())
		return err
	}
	// Setup the specified GPIO pin
	pin := rpio.Pin(pinNumber)
	// Specify that the pin should be an output
	pin.Output()
	// Handle the power change
	if powerOn {
		pin.High()
	} else {
		pin.Low()
	}
	// Unmap /dev/mem when done
	return rpio.Close()
}
