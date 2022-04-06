package firmware

import (
	"errors"
	"fmt"

	"github.com/MikMuellerDev/smarthome-hw/core/config"
	"github.com/MikMuellerDev/smarthome-hw/core/log"
)

// Acts as a global lock, in this case if a code is being sent in order to prevent interrupts
var blocked bool

var (
	ErrNoCode   = errors.New("can not process request: switch has no code entry")
	ErrBlocked  = errors.New("can not process request now: sender busy")
	ErrDisabled = errors.New("hardware is currently disabled")
)

// If the sender's hardware is not blocked, the code corresponding to the function's input is sent
// Returns an error if the sender is uninitialized or the switch has no entry in the config file
func SetPower(switchId string, powerOn bool) error {
	if blocked {
		log.Warn("Can not send code right now: the sender is currently busy")
		return ErrBlocked
	}
	if !config.GetConfig().Hardware.HardwareEnabled {
		log.Warn("Can not send code right now: the hardware is currently disabled")
		return ErrDisabled
	}
	blocked = true
	if err := Init(); err != nil {
		log.Error("Initializing sender failed: ", err.Error())
		blocked = false
		return err
	}
	switches := config.GetConfig().Switches
	for _, switchItem := range switches {
		if switchItem.Id == switchId {
			var code int
			if powerOn {
				code = switchItem.CodeOn
			} else {
				code = switchItem.CodeOff
			}
			if err := sender.Send(code); err != nil {
				log.Error("Failed to send code: ", err.Error())
				blocked = false
				return err
			}
			log.Info(fmt.Sprintf("Successfully send code %d (Switch: %s PowerOn: %t)", code, switchId, powerOn))
			if err := Free(); err != nil {
				log.Error("Deactivated sender after transmit failed: ", err.Error())
				blocked = false
				return err
			}
			blocked = false
			return nil
		}
	}
	blocked = false
	return ErrNoCode
}

// Sends a code without any preproccessing
func SendCode(code int) error {
	if blocked {
		log.Warn("Can not send code right now: the sender is currently busy")
		return ErrBlocked
	}
	if !config.GetConfig().Hardware.HardwareEnabled {
		log.Warn("Can not send code right now: the hardware is currently disabled")
		return ErrDisabled
	}
	blocked = true
	if err := Init(); err != nil {
		log.Error("Initializing sender failed: ", err.Error())
		blocked = false
		return err
	}
	if err := sender.Send(code); err != nil {
		log.Error("Failed to send code: ", err.Error())
		blocked = false
		return err
	}
	if err := Free(); err != nil {
		log.Error("Deactivated sender failed: ", err.Error())
		blocked = false
		return err
	}
	blocked = false
	return nil
}
