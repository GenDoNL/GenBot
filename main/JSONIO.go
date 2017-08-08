package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Reads the server data from the data file.
func readServerData() {
	raw, err := ioutil.ReadFile(DataLocation)
	if err != nil {
		fmt.Println(err.Error())
		if len(Servers) == 0 {
			Servers = make(map[string]*(ServerData))
		}

	} else {
		json.Unmarshal(raw, &Servers)
	}

}

// Writes the server data to the data file.
func writeServerData() {

	data, err := json.MarshalIndent(Servers, "", "    ")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rc := ioutil.WriteFile(DataLocation, data, 0644)
	if rc != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
