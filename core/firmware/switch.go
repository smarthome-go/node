package firmware

import (
	"errors"
	"fmt"
	"sync"

	"github.com/smarthome-go/node/core/config"
	"github.com/smarthome-go/node/core/log"
)

// Acts as a global lock, in this case if a code is being sent in order to prevent interrupts
var blocked sync.Mutex

var (
	ErrNoSwitch = errors.New("can not process request: switch has no RF-code or GPIO-pin entry")
	ErrBlocked  = errors.New("can not process request now: hardware busy")
	ErrDisabled = errors.New("hardware is currently disabled")
)

// If the sender's hardware is not blocked, the code corresponding to the function's input is sent
// Returns an error if the sender is uninitialized or the switch has no entry in the config file
func SetPower(
	switchId string,
	powerOn bool,
) error {
	blocked.Lock()
	defer blocked.Unlock()

	if !config.GetConfig().Hardware.HardwareEnabled {
		log.Trace("Can not send code right now: the hardware is currently disabled")
		return ErrDisabled
	}

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
			if err := sendCode(
				code,
				config.GetConfig().Hardware,
			); err != nil {
				log.Error("Failed to send code: ", err.Error())
				return err
			}
			log.Trace(fmt.Sprintf("Successfully send code `%d`. (Switch: %s | PowerOn: %t)", code, switchId, powerOn))
			return nil
		}
	}

	// Handle GPIO switches afterwards
	switchesGPIO := config.GetConfig().SwitchesGPIO
	for _, switchItem := range switchesGPIO {
		if switchItem.Id == switchId {
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
	// Check if the missing switch is ignored
	switchesIgnore := config.GetConfig().SwitchesIgnore
	for _, switchItem := range switchesIgnore {
		if switchItem == switchId {
			log.Debug(fmt.Sprintf("Skipping switch `%s`: switch is ignored", switchId))
			return nil
		}
	}
	log.Warn(fmt.Sprintf("Unregistered switch requested: switch `%s` was requested but not found.", switchId))
	return ErrNoSwitch
}

// Sends a code without any preproccessing
func SendCode(code int) error {
	blocked.Lock()
	defer blocked.Unlock()

	if !config.GetConfig().Hardware.HardwareEnabled {
		log.Trace("Can not send code right now: the hardware is currently disabled")
		return ErrDisabled
	}
	if err := sendCode(
		code,
		config.GetConfig().Hardware,
	); err != nil {
		log.Error("Failed to send code: ", err.Error())
		return err
	}

	return nil
}
