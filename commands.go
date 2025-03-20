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

	err = s.db.DeleteAllFeeds(ctx)
	if err != nil {
		return err
	}
	log.Println("Successfully deleted all feeds from feeds table!")

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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		err := errors.New("You must supply both a name for the feed and the url, in that order")
		return err
	}
	// Need to get the current user's id
	ctx := context.Background()
	// For possible null uuid
	nullUUID := uuid.NullUUID{
		UUID:  user.ID,
		Valid: true,
	}

	// Make new feed
	newID := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt
	feedName := cmd.args[0]
	feedURL := cmd.args[1]
	feedArgs := database.CreateFeedParams{
		ID:        newID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      feedName,
		Url:       feedURL,
		UserID:    nullUUID,
	}

	feed, err := s.db.CreateFeed(ctx, feedArgs)
	if err != nil {
		return err
	}

	// If everything went well:
	// we want to print out the fields of the new feed record
	log.Printf("%#v\n", feed)

	// We also want to make a new feed follow for the current user
	newCmd := command{
		args: []string{feedArgs.Url},
		name: "follow",
	}
	handlerFollow(s, newCmd, user)

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	ctx := context.Background()
	allFeeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, feed := range allFeeds {
		row := fmt.Sprintf("Feed Name: %v\tUrl: %v\tAdded by: %v\n", feed.Feedname, feed.Feedurl, feed.Username)
		log.Println(row)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	// If no url arg, fail
	if len(cmd.args) < 1 {
		newErr := errors.New("No url provided. Please add a url to be able to follow a new feed")
		return newErr
	}
	url := cmd.args[0]
	// Add context
	ctx := context.Background()

	// Set up the feed follow params
	//user, err := s.db.GetUser(ctx, (*s).cfg.CurrentUserName)
	//if err != nil {
	//	return err
	//}
	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}

	valueToInsert := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
		FeedID: uuid.NullUUID{
			UUID:  feed.ID,
			Valid: true,
		},
	}
	// Create a new feed follow record for the current user
	createdFeedFollow, err := s.db.CreateFeedFollow(ctx, valueToInsert)
	if err != nil {
		return err
	}
	// Should print the name of the feed
	// and the current user once the record is created
	newFeedFollow := fmt.Sprintf("Feed Name: %v\nCurrent user: %v\n", createdFeedFollow.FeedName, createdFeedFollow.UserName)
	log.Println(newFeedFollow)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	// Add context
	ctx := context.Background()
	// Get current user's id
	// user, err := s.db.GetUser(ctx, (*s).cfg.CurrentUserName)
	// if err != nil {
	// 	return err
	// }
	// Get all feeds followed by user
	feedsFollowed, err := s.db.GetFeedFollow(ctx, uuid.NullUUID{
		UUID:  user.ID,
		Valid: true,
	})
	if err != nil {
		return err
	}
	// Print out all feedsFollowed
	log.Printf("%v is following:\n", user.Name)
	for _, feed := range feedsFollowed {
		row := fmt.Sprintf("* %v\n", feed)
		log.Println(row)
	}

	return nil
}
