# gator

A terminal RSS aggregator that allows you to follow feeds and browse the latest posts.

This is a guided project from the [Boot.dev](https://boot.dev) curriculum.

The project is built in Go using PostgreSQL, sqlc for type-safe SQL queries, and goose for schema migrations. 

## Features

- **Follow feeds and browse the latest posts.**
- **Add users and switch between them.**
- **Follow feeds that other users have added.**
- **Background worker for periodic RSS feed scraping.** 

## Prerequisites

Make sure you have the following installed on your system:

- **Go** (version 1.27.1 or higher)
- **PostgreSQL** (running locally)
- **goose** (required for running database migrations)

## Installation & Setup

### 1. Clone the repository

```bash
git clone https://github.com/lbergold-dev/gator.git
cd gator
```

### 2. Configure the application

Create the configuration file at `~/.config/gator/config.json`:

```json
{
  "db_url": "postgres://localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

### 3. Run database migrations

Start PostgreSQL, then create the database:

```bash
brew services start postgresql@15
psql postgres -c "CREATE DATABASE gator;"
```

Apply the SQL schema migrations using `goose`:

```bash
goose -dir sql/schema postgres "postgres://localhost:5432/gator?sslmode=disable" up
```

### 4. Build or install the binary

Install the CLI globally into your `$GOBIN`:

```bash
go install .
```

Or build it locally:

```bash
go build -o gator .
```

## Usage

```bash
# Register a new user
gator register alice

# Log in as an existing user
gator login alice

# Add and automatically follow a feed
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"

# Start the background aggregator (polls feeds every 1 minute)
gator agg 1m

# Browse the latest posts (default limit: 2)
gator browse 5
```

### Available Commands

| Command | Arguments | Description |
| :--- | :--- | :--- |
| `register` | `<username>` | Creates a new user and sets them as active |
| `login` | `<username>` | Switches the active user |
| `users` | — | Lists all registered users |
| `addfeed` | `<name> <url>` | Adds a new feed and follows it for the current user |
| `feeds` | — | Lists all registered feeds and their creators |
| `follow` | `<url>` | Follows an existing feed |
| `following` | — | Lists all feeds followed by the current user |
| `unfollow` | `<url>` | Unfollows a feed for the current user |
| `agg` | `<duration>` | Runs the background feed scraper (e.g., `30s`, `1m`) |
| `browse` | `[limit]` | Displays the latest posts from followed feeds |
| `reset` | — | Deletes all users and cascades related data |

## Architecture & Project Structure

```text
.
├── internal/
│   ├── config/       # Local JSON configuration reader/writer
│   └── database/     # sqlc-generated type-safe database layer
├── sql/
│   ├── queries/      # Raw SQL queries for sqlc
│   └── schema/       # Goose database migrations
├── commands.go       # CLI command registry and dispatcher
├── middleware.go     # Authentication middleware for user-scoped commands
├── rss.go            # HTTP client, XML stream decoder, and feed scraper
├── handler_*.go      # Command handlers
└── main.go           # Application entry point and state initialization
```

## Development

To regenerate the Go database code after modifying files in `sql/queries/` or `sql/schema/`:

```bash
sqlc generate
```

## License

MIT -- see [LICENSE](LICENSE)
