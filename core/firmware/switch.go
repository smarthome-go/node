package firmware

import (
	"errors"
	"fmt"

	"github.com/smarthome-go/node/core/config"
	"github.com/smarthome-go/node/core/log"
)

// Acts as a global lock, in this case if a code is being sent in order to prevent interrupts
var blocked bool

var (
	ErrNoSwitch = errors.New("can not process request: switch has no RF-code or GPIO-pin entry")
	ErrBlocked  = errors.New("can not process request now: hardware busy")
	ErrDisabled = errors.New("hardware is currently disabled")
)

// If the sender's hardware is not blocked, the code corresponding to the function's input is sent
// Returns an error if the sender is uninitialized or the switch has no entry in the config file
func SetPower(switchId string, powerOn bool) error {
	if blocked {
		log.Trace("Can not send code right now: the sender is currently busy")
		return ErrBlocked
	}
	if !config.GetConfig().Hardware.HardwareEnabled {
		log.Trace("Can not send code right now: the hardware is currently disabled")
		return ErrDisabled
	}
	blocked = true
	// Handle RF switchesRF first
	switchesRF := config.GetConfig().SwitchesRF
	for _, switchItem := range switchesRF {
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
			log.Trace(fmt.Sprintf("Successfully send code `%d`. (Switch: %s | PowerOn: %t)", code, switchId, powerOn))
			blocked = false
			return nil
		}
	}
	// Handle GPIO switches afterwards
	switchesGPIO := config.GetConfig().SwitchesGPIO
	for _, switchItem := range switchesGPIO {
		if switchItem.Id == switchId {
			blocked = false
			log.Trace(fmt.Sprintf("Successfully handled switch GPIO. (Switch: %s | PowerOn: %t) ", switchId, powerOn))
			// If the switch uses the `invert` setting, invert the power state
			if switchItem.Invert {
				return setPin(
					switchItem.Pin,
					!powerOn,
				)
			}
			// Otherwise, use the normal setting
			return setPin(
				switchItem.Pin,
				powerOn,
			)
		}
	}
	log.Warn(fmt.Sprintf("Unregistered switch requested: switch `%s` was requested but not found.", switchId))
	blocked = false
	return ErrNoSwitch
}

// Sends a code without any preproccessing
func SendCode(code int) error {
	if blocked {
		log.Trace("Can not send code right now: the sender is currently busy")
		return ErrBlocked
	}
	if !config.GetConfig().Hardware.HardwareEnabled {
		log.Trace("Can not send code right now: the hardware is currently disabled")
		return ErrDisabled
	}
	blocked = true
	if err := sender.Send(code); err != nil {
		log.Error("Failed to send code: ", err.Error())
		blocked = false
		return err
	}
	blocked = false
	return nil
}
