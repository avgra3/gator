package main

import (
	"fmt"
	"internal/config"
	"log"
	//"github.com/avgra3/gator/internal/config"
)

func main() {
	// read the config file
	currentConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	// set current user to lane
	// save config to disk
	err = config.SetUser("lane", &currentConfig)
	if err != nil {
		log.Fatal(err)
	}
	// read config and print contents of config struct to terminal
	readConfigFile, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	output := fmt.Sprintf("DB Url: %v\nCurrent Username: %v", readConfigFile.DBUrl, readConfigFile.CurrentUserName)
	fmt.Println(output)
}
