package main

import (
	"encoding/json"
	"io/ioutil"
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

// Disabled until fully ported to database
//// Reads the server data from the data file.
//func readServerData() {
//	raw, err := ioutil.ReadFile(BotConfig.DataLocation)
//	if err != nil {
//		log.Error(err.Error())
//		if len(Servers) == 0 {
//			Servers = make(map[string]*ServerData)
//			log.Infof("Initialized server data.")
//		}
//	} else {
//		json.Unmarshal(raw, &Servers)
//		log.Debugf("Read data from file at %s", BotConfig.DataLocation)
//	}
//
//}
//
//func writeSeperateFiles() {
//	for _, server := range Servers {
//		data, err := json.MarshalIndent(server, "", "    ")
//		if err != nil {
//			fmt.Println(err.Error())
//			os.Exit(1)
//		}
//
//
//		location := BotConfig.DataLocation[:len(BotConfig.DataLocation)-9] + "temp/" + server.ID
//		rc := ioutil.WriteFile(location, data, 0644)
//		if rc != nil {
//			log.Error("Issue encountered while writing to file, shutting down.")
//			log.Error(location)
//			log.Fatal(rc.Error())
//		}
//	}
//}
