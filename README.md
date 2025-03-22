# 🐊 Gator 🐊

A simple cli RSS aggregator written in Go using PostgreSQL.

## Tools Needed to Run/Build

In order to have the app work as expected, you will need the following running.

- [PostgreSQL](https://www.postgresql.org/)
    - Once your Postgres server is setup, you will need to ensure there is a user with name: "postgres"
- [sqlc](https://sqlc.dev/) _only needed if you are planning on extending functionality to gator._
- [goose](https://github.com/pressly/goose)
- [go](https://go.dev/)
- Copy [.gatorconfig.json](./.gatorconfig) to your home directory.
    - On Linux this would be something like "/home/youUserName/.gatorconfig.json"
    - Update it with the below:

```json
{"db_url":"postgres://USERNAME:PASSWORD@IP_ADDRESS:PORT/gator?sslmode=disable","current_user_name":"madeUpName"}
```
- USERNAME: whatever default username you set upon postgresql setup
- PASSWORD: associated password
- IP_ADDRESS: `localhost` if the database is on the same machine, otherwise, the database ip address
- PORT: `5432` is the default

## Build and install from source

Below is an example of building on Linux

```bash
git clone --depth 1 https://github.com/avgra3/gator

cd gator

goose -dir ./sql/schema up

go install
```

## Usage

```bash
gator COMMAND ARGS
```

### Available Commands

- "login" USERNAME => Will fail if the user doesn't exists.
- "register" USERNAME => Will fail if the user already exists.
- "reset" => Resets the database. Useful for testing.
- "users" => Returns all users with indication of the current user.
- "agg" TIME_DURATION => with a given time duration, will go through all feeds and pull posts, starting from the oldest upate.
- "addfeed" FEED_NAME FEED_URL => Add a new feed to feeds. This will automatically follow the feed for the current user.
- "feeds" => Returns all feeds and who added them (which does not necessarily mean the user who added the feed.
- "follow" FEED_URL => Add a url to have the user follow
- "following" => See who the current logged in user is following.
- "unfollow" FEED_URL => You must provide a url to unfollow.
- "browse" LIMIT (optional, default 2) => Show posts with title name, description, and url to follow to the actual article.

