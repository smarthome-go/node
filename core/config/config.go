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

type Config struct {
	Port      uint16   `json:"port"`
	NodeName  string   `json:"nodeName"`
	TokenHash string   `json:"tokenHash"`
	Hardware  Hardware `json:"hardware"`
	Switches  []Switch `json:"switches"`
}

var Version string

// Documentation of following parameters: github.com/smarthome-go/rpirf
type Hardware struct {
	HardwareEnabled     bool   `json:"hardwareEnabled"`
	RFDevicePin         uint8  `json:"pin"` // The BCM pin to which a 433mhz sender is attached
	RFDeviceProtocol    uint8  `json:"protocol"`
	RFDeviceRepeat      uint8  `json:"repeat"`
	RFDevicePulselength uint16 `json:"pulseLength"`
	RFDeviceLength      uint8  `json:"contentLength"`
}

type Switch struct {
	Id      string `json:"id"`
	CodeOn  int    `json:"on"`
	CodeOff int    `json:"off"`
}

var config Config

// The path were the config file is located
const configPath = "./config.json"

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
			log.Fatal("Failed to initialize config: could not read or create a config file: ", errCreate.Error())
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
			HardwareEnabled:     false,
			RFDevicePin:         0,
			RFDeviceProtocol:    1,
			RFDeviceRepeat:      10,
			RFDevicePulselength: 180,
			RFDeviceLength:      24,
		},
		Switches: []Switch{
			{
				Id:      "s1",
				CodeOn:  0,
				CodeOff: 0,
			},
		},
	}
	fileContent, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		log.Error("Failed to create config file: creating file content from JSON failed: ", err.Error())
		return Config{}, err
	}
	if err = ioutil.WriteFile("./config.json", fileContent, 0644); err != nil {
		log.Error("Failed to write file to disk: ", err.Error())
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
		log.Fatal("Error during unmarshal: ", err.Error())
		return err
	}
	configJson, _ := json.MarshalIndent(&config, "", "    ")
	err = ioutil.WriteFile("./config.json", configJson, 0644)
	if err != nil {
		log.Fatal("Error writing new token hash to config.json: ", err.Error())
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
