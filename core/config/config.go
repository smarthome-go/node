package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/MikMuellerDev/smarthome-hw/core/log"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Port       uint16
	Production bool
	NodeName   string
	TokenHash  string
	Version    string
	GOARCH     string
}

// Documentation of following parameters: github.com/MikMuellerDev/rpirf
type Hardware struct {
	RFDevicePin         uint8 // The BCM pin to which a 433mhz sender is attached
	RFDeviceProtocol    uint8
	RFDeviceRepeat      uint8
	RFDevicePulselength uint16
	RFDeviceLength      uint8
}

type Switch struct {
	Id      string
	CodeOn  int
	CodeOff int
}

var config Config

// The path were the config file is located
const configPath = "./config.json"

func GenerateTokenHash(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to generate new hash: ", err.Error())
		return "", err
	}
	return string(hash), nil
}

func ReadConfigFile() error {
	// Read file from <configPath> on disk
	// If this file does not exist, create a new blank one
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
	configFile.GOARCH = runtime.GOARCH
	config = configFile
	return nil
}

// Creates an empty config file, can return an error if it fails
func createNewConfigFile() (Config, error) {
	hash, err := GenerateTokenHash("smarthome")
	if err != nil {
		return Config{}, err
	}
	fmt.Printf("The default authentication token is\x1b[1;34m 'smarthome' \x1b[1;0m\nIt is recommended to change it\x1b[1;31m immediately\x1b[1;0m (using update_token.sh).\n")
	config := Config{
		Port:       8081,
		Production: false,
		NodeName:   "localhost",
		TokenHash:  hash, // TODO: generate with bcrypt
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

func GetConfig() Config {
	return config
}

func SetVersion(version string) {
	config.Version = version
}

func SetHash(hash string) {
	config.TokenHash = hash
}
