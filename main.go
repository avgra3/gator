package main

import (
	// "fmt"
	"internal/config"
	"log"
	"os"
	//"github.com/avgra3/gator/internal/config"
)

func main() {
	// read the config file
	currentConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	// Create a new state
	s := state{ptrConfig: &currentConfig}

	// Create commands struct
	cmds := commands{
		commandNames: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	// Getting current args
	args := os.Args
	// Error if < 2 args
	if len(args) < 2 {
		// The first arg is the program name
		log.Fatal("There needs to be at least 1 argument present")
	}
	// Need to split command-line args intwo command name and args slice
	// to create a command instance
	cmdName := args[1]
	cmdArgs := args[2:]

	// use commands.run method to run the given command
	cmd := command{
		args: cmdArgs,
		name: cmdName,
	}
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
	// Write config to file
	config.SetUser(cmdArgs[0], &currentConfig)

	// read config and print contents of config struct to terminal
	// readConfigFile, err := config.Read()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// output := fmt.Sprintf("DB Url: %v\nCurrent Username: %v", readConfigFile.DBUrl, readConfigFile.CurrentUserName)
	// fmt.Println(output)
}
