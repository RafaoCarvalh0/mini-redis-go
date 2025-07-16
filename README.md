# mini-redis-go

This project is a practical study to get familiar with the Go programming language and to understand the internal workings of a minimalist Redis server. The main goal is to explore concepts such as networking, binary protocols (RESP2), file handling, and data structures in Go, as well as to learn about persistence and client-server communication.

## Motivation

- **Learn Go**: Practice the Go language in a real TCP server context.
- **Understand Redis**: Learn how an in-memory database works, including basic commands and persistence.
- **Explore protocols**: Implement the RESP2 protocol, used by Redis for client communication.
- **Persistence**: Simulate data snapshots on disk, similar to Redis RDB.

## Project structure

```
app/
  commands/           # Core Redis command handlers (SET, GET, etc)
  protocol_parser/    # RESP2 protocol parser
  server_config/      # Server configuration logic
  rdb_utils.go        # RDB file handling utilities
  main.go             # Server entry point

cli/
  mini_redis_go_cli.go # CLI client to send commands to the server

README.md             # Project documentation and usage
go.mod, go.sum        # Go module dependencies
.env                  # (Optional) Environment variables for host/port
template_dump.rdb     # Example RDB file for persistence
```

## Requirements

- Go 1.18+

## Setup

1. **Clone the repository**
2. **Install dependencies:**
   ```sh
   go mod tidy
   ```
3. **(Optional) Create a `.env` file** in the project root to configure server/CLI defaults:
   ```env
   MINI_REDIS_HOST=127.0.0.1
   MINI_REDIS_PORT=6379
   ```
   If not set, defaults are used.

## Running the server

From the `app` directory:
```sh
cd app
# Run with Go
go run .
```

You can customize host, port, data directory, and db filename with flags or environment variables:
```sh
go run . --host 0.0.0.0 --port 6379 --dir "../" --dbfilename "template_dump.rdb"
```

## Running the CLI

The CLI binary (`mini-redis-cli.exe`) is already built and available in the project root. You can use it directly to send commands:

```sh
./mini-redis-cli SET foo bar
./mini-redis-cli GET foo
./mini-redis-cli SET temp value PX 5000
```

You can also specify host and port:
```sh
./mini-redis-cli --host 127.0.0.1 --port 6379 GET foo
```

> **⚠️ Important:**
> 
> The CLI **must use the same host and port as the server**. If you started the server with custom `--host` or `--port` values, use the same values when running the CLI. Otherwise, the connection will fail.

## Supported commands

- `PING`                   - Test the connection with the server
- `ECHO <message>`         - Echo back the provided message
- `SET <key> <value> [PX milliseconds]` - Set a value for a key, optionally with expiration in ms (PX)
- `GET <key>`              - Get the value of a key
- `CONFIG <subcommand>`    - Manage server configuration
- `KEYS <pattern>`         - List keys matching the pattern
- `SAVE`                   - Save the current dataset to disk

## Notes

- This project is **not recommended for production use**.
- The focus is on learning, not on performance or full Redis compatibility.
- The implementation covers only basic commands and simplified persistence.

## Credits

Developed as part of a personal study and inspired by the [Codecrafters Redis Go challenge](https://codecrafters.io/challenges/redis). 