package main

import (
	"database/sql"
	"github.com/avgra3/gator/internal/database"
	_ "github.com/lib/pq"
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
	// DB Connection
	db, err := sql.Open("postgres", currentConfig.DBUrl)
	dbQueries := database.New(db)

	// Create a new state
	s := state{cfg: &currentConfig, db: dbQueries}

	// Create commands struct
	cmds := commands{
		commandNames: make(map[string]func(*state, command) error),
	}
	// Registered commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)

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
}
