package main

import (
	"context"

	"github.com/avgra3/gator/internal/database"
)

// Function to allow us to verify we are logged in
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	// Return a new function that matches the standard handler signature
	return func(s *state, cmd command) error {
		// 1. Check if user is logged in by retrieving current user
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		// 2. Since the user is correct, call original handler with user
		return handler(s, cmd, user)
	}
}
