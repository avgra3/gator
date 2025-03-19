package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/avgra3/gator/internal/database"
	"github.com/google/uuid"
	"internal/config"
	"log"
	"time"
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

	// Save the config file
	config.SetUser((*s).cfg.CurrentUserName, (*s).cfg)

	// Send message to cli
	message := fmt.Sprintf("The user \"%v\" has been set", (*s).cfg.CurrentUserName)
	log.Println(message)
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
		return err
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

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		return err
	}
	log.Println("Successfully deleted all users from users table!")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	ctx := context.Background()
	names, err := s.db.GetUsers(ctx)
	currentUser := s.cfg.CurrentUserName
	if err != nil {
		return err
	}
	for _, name := range names {
		if name == currentUser {
			log.Printf("* %v (current)\n", name)
		} else {
			log.Printf("* %v\n", name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	baseUrl := "https://www.wagslane.dev/index.xml"
	rssFeedPtr, err := fetchFeed(ctx, baseUrl)
	if err != nil {
		return err
	}
	log.Printf("%#v\n", rssFeedPtr)
	return nil
}
