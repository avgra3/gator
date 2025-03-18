package config

import (
	"errors"
	"fmt"
)

// Types needed
type state struct {
	ptrConfig *Config
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
		err := errors.New("The login handler expects a single argument: ther username, which was not provided")
		return err
	}

	(*s).ptrConfig.CurrentUserName = cmd.args[0]
	message := fmt.Sprintf("The user \"%v\" has been set", (*s).ptrConfig.CurrentUserName)
	fmt.Println(message)
	return nil
}
