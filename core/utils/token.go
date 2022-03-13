package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/MikMuellerDev/smarthome-hw/core/config"
	"github.com/MikMuellerDev/smarthome-hw/core/log"
	"golang.org/x/crypto/bcrypt"
)

func WriteHashToConfig(token string) error {

	var jsonBlob = []byte(`{}`)
	config := config.GetConfig()
	err = json.Unmarshal(jsonBlob, &config)
	if err != nil {
		log.Fatal("Error during unmarshal: ", err.Error())
		return err
	}
	configJson, _ := json.MarshalIndent(&config, "", "    ")
	err = ioutil.WriteFile("./config.json", configJson, 0644)
	if err != nil {
		log.Fatal("Error writing  new token hash to config.json: ", err.Error())
		return err
	}
	log.Debug("Written new token hash to config.json")
	return nil
}
