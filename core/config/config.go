package config

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"

	"golang.org/x/crypto/bcrypt"

	"github.com/smarthome-go/node/core/log"
)

var Version string

// The path were the config file is located
const configPath = "./config.json"

var config Config

type Config struct {
	Port           uint16       `json:"port"`
	NodeName       string       `json:"nodeName"`
	TokenHash      string       `json:"tokenHash"`
	Hardware       Hardware     `json:"hardware"`
	SwitchesRF     []SwitchRF   `json:"switchesRF"`
	SwitchesGPIO   []SwitchGPIO `json:"switchesGPIO"`
	SwitchesIgnore []string     `json:"switchesIgnore"`
}

// Documentation of following parameters: github.com/smarthome-go/rpirf
type Hardware struct {
	UseGpioCdev         bool    `json:"useGpioCdev"`
	GpioCdevChip        *string `json:"gpioCdevChip"`
	HardwareEnabled     bool    `json:"hardwareEnabled"`
	RFDevicePin         int     `json:"pin"` // The BCM pin to which a 433mhz sender is attached
	RFDeviceProtocol    uint8   `json:"protocol"`
	RFDeviceRepeat      uint8   `json:"repeat"`
	RFDevicePulselength uint16  `json:"pulseLength"`
	RFDeviceLength      uint8   `json:"contentLength"`
}

type SwitchRF struct {
	Id      string `json:"id"`  // Id used by Smarthome server
	CodeOn  int    `json:"on"`  // Code to send when the Smarthome server requests power ON
	CodeOff int    `json:"off"` // Code to send when the Smarthome server requests power OFF
}

type SwitchGPIO struct {
	Id     string `json:"id"`     // Id used by Smarthome server
	Pin    uint8  `json:"pin"`    // The BCM pin to which a GPIO device is attached
	Invert bool   `json:"invert"` // Whether the power request should be treated in an inverted manner (mainly useful for relay control)
}

// A dry-run of the `RadConfigFile()` method used in the healthtest
func ProbeConfigFile() error {
	// Read file from <configPath> on disk
	// If this file does not exist, return an error
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Error("Failed to read config file: ", err.Error())
		return nil
	}
	// Parse config file to a test struct <Config>
	var configFile Config
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&configFile)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to parse config file at `%s` into Config struct: %s", configPath, err.Error()))
		return err
	}
	return nil
}

// Reads the config file from disk, if the file does not exist (for example first run), a new one is created
func ReadConfigFile() error {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		configTemp, errCreate := createNewConfigFile()
		if errCreate != nil {
			log.Error("Failed to read config file: ", err.Error())
			log.Error("Failed to initialize config: could not read or create a config file: ", errCreate.Error())
			return err
		}
		config = configTemp
		log.Info("Failed to read config file: but managed to create a new config file")
		return nil
	}
	// Parse config file to struct <Config>
	var configFile Config
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&configFile)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to parse config file at `%s` into Config struct: %s", configPath, err.Error()))
		return err
	}

	// validate the semantic correctness of the config file
	if configFile.Hardware.UseGpioCdev == (configFile.Hardware.GpioCdevChip == nil) {
		return fmt.Errorf("Invalid CDEV configuration: option `use cdev` must correspond to value `cdev chip`")
	}

	config = configFile
	return nil
}

// Creates an empty config file, can return an error if it fails
func createNewConfigFile() (Config, error) {
	token := randomBase64String(32)
	hash, err := GenerateTokenHash(token)
	if err != nil {
		return Config{}, err
	}
	log.Info("The default authentication token was generated: it is printed in the logs below")
	fmt.Printf("INFO[0000] The token is\x1b[1;34m '%s' \x1b[1;0m\n", token)
	log.Info("It can be changed using update_token.sh")
	config := Config{
		Port:      8081,
		NodeName:  "localhost",
		TokenHash: hash,
		Hardware: Hardware{
			UseGpioCdev:         false,
			GpioCdevChip:        nil,
			HardwareEnabled:     false,
			RFDevicePin:         0,
			RFDeviceProtocol:    1,
			RFDeviceRepeat:      10,
			RFDevicePulselength: 180,
			RFDeviceLength:      24,
		},
		SwitchesRF: []SwitchRF{
			{
				Id:      "s1",
				CodeOn:  0,
				CodeOff: 0,
			},
		},
		SwitchesGPIO: []SwitchGPIO{
			{
				Id:  "s2",
				Pin: 1,
			},
		},
		SwitchesIgnore: []string{
			"s3",
		},
	}
	fileContent, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		log.Error("Failed to create config file: creating file content from JSON failed: ", err.Error())
		return Config{}, err
	}
	if err = ioutil.WriteFile("./config.json", fileContent, 0644); err != nil {
		log.Error("Failed to write config file to disk: ", err.Error())
		return Config{}, err
	}
	return config, nil
}

// Writes the current state of the global config to the file on disk
func WriteConfig() error {
	var jsonBlob = []byte(`{}`)
	config := config
	err := json.Unmarshal(jsonBlob, &config)
	if err != nil {
		log.Error("Error during unmarshal: ", err.Error())
		return err
	}
	configJson, _ := json.MarshalIndent(&config, "", "    ")
	err = ioutil.WriteFile("./config.json", configJson, 0644)
	if err != nil {
		log.Error("Error writing new token hash to config.json: ", err.Error())
		return err
	}
	log.Debug("Written to config.json")
	return nil
}

/*
Set / Get functions
*/

func GetConfig() Config {
	return config
}

func SetHash(hash string) {
	config.TokenHash = hash
}

func SetHardwareEnabled(enabled bool) {
	config.Hardware.HardwareEnabled = enabled
}

/*
Token-Related
*/

// Returns a random string of a given length that is encoded in BS64
func randomBase64String(length int) string {
	buff := make([]byte, int(math.Ceil(float64(length)/float64(1.33333333333))))
	if _, err := rand.Read(buff); err != nil {
		log.Error("Encoding json failed: ", err.Error())
	}
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:length] // strip 1 extra character we get from odd length results
}

func GenerateTokenHash(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), 5)
	if err != nil {
		log.Error("Failed to generate new hash: ", err.Error())
		return "", err
	}
	return string(hash), nil
}
