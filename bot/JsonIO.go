package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Reads the config file.
func readConfig() {
	raw, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		json.Unmarshal(raw, &BotConfig)
		log.Info("Config file loaded.")
	}
}

// Reads the server data from the data file.
func readServerData() {
	raw, err := ioutil.ReadFile(BotConfig.DataLocation)
	if err != nil {
		log.Error(err.Error())
		if len(Servers) == 0 {
			Servers = make(map[string]*ServerData)
			log.Infof("Initialized server data.")
			writeServerData()
		}
	} else {
		json.Unmarshal(raw, &Servers)
		log.Debugf("Read data from file at %s", BotConfig.DataLocation)
	}

}

// Writes the server data to the data file.
func writeServerData() {
	data, err := json.MarshalIndent(Servers, "", "    ")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rc := ioutil.WriteFile(BotConfig.DataLocation, data, 0644)
	if rc != nil {
		log.Error("Issue encountered while writing to file, shutting down.")
		log.Fatal(err.Error())
	}
}
