package main

import (
	"context"
	"errors"
	"fmt"
	"internal/config"
	"log"
	"time"

	"github.com/avgra3/gator/internal/database"
	"github.com/google/uuid"
)

// Types needed
type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	args []string
	name string
}

type commands struct {
	commandNames map[string]func(*state, command) error
}

// Registers a new handler function for a command name
func (c *commands) register(name string, f func(*state, command) error) {
	c.commandNames[name] = f
}

// Runs a given command with the provided state if it exists
func (c *commands) run(s *state, cmd command) error {
	name := cmd.name
	err := c.commandNames[name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		err := errors.New("The login handler expects a single argument: the username, which was not provided")
		return err
	}
	// Error if user does not exist:
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, cmd.args[0])
	if err != nil {
		log.Println("here")
		return err
	}
	if (user.Name == "" || user == database.User{}) {
		err = errors.New("Unknown user")
		return err
	}
	(*s).cfg.CurrentUserName = cmd.args[0]
	message := fmt.Sprintf("The user \"%v\" has been set", (*s).cfg.CurrentUserName)

	fmt.Println(message)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		err := errors.New("A username must be provided in order for a user to be registered.")
		return err
	}
	(*s).cfg.CurrentUserName = cmd.args[0]
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	ctx := context.Background()
	user, err := s.db.CreateUser(ctx, newUser)
	if err != nil {
		log.Fatal(err)
	}
	// Success message
	log.Printf("User \"%v\" successfully added!\n", newUser.Name)
	// Set user to current user
	newCmd := command{
		args: []string{user.Name},
		name: "login",
	}
	handlerLogin(s, newCmd)

	return nil
}
