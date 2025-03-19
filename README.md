# Gator

An RSS aggregator written in Go using PostgreSQL for the database.

## Tools Used

- PostgreSQL
- sqlc
- goose

## Available Commands

- "login" USERNAME => Will fail if the user doesn't exists.
- "register" USERNAME => Will fail if the user already exists.
- "reset" => Resets the database. Useful for testing.
- "users" => Returns all users with indication of the current user.
