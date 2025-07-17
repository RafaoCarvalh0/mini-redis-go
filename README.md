# mini-redis-go

This project is a practical study to get familiar with the Go programming language and to understand the internal workings of a minimalist Redis server. The main goal is to explore concepts such as networking, binary protocols (RESP2), file handling, and data structures in Go, as well as to learn about persistence and client-server communication.

**Note:** This implementation currently saves and loads only string values.

## Motivation

- **Learn Go**: Practice the Go language in a real TCP server context.
- **Understand Redis**: Learn how an in-memory database works, including basic commands and persistence.
- **Explore protocols**: Implement the RESP2 protocol, used by Redis for client communication.
- **Persistence**: Simulate data snapshots on disk, similar to Redis RDB.

## Project structure

```
app/
  commands/           # Core Redis command handlers (SET, GET, DEL, etc)
  protocol_parser/    # RESP2 protocol parser
  server_config/      # Server configuration logic
  rdb_utils.go        # RDB file handling utilities
  rdb_constants.go    # RDB format constants
  main.go             # Server entry point

cli/
  mini_redis_go_cli.go # CLI client to send commands to the server

example_dump.rdb      # Example RDB file for persistence (gitignored)
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

## Compiling the CLI

From the project root, run:

**On Windows:**
```sh
go build -o mini-redis-cli.exe cli/mini_redis_go_cli.go
```

**On Linux/Mac:**
```sh
go build -o mini-redis-cli cli/mini_redis_go_cli.go
```

This will generate the executable in the project root. You can then use it as described in the section below.

## Running the CLI

On Linux:
```sh
./mini-redis-cli help
```

On windows:
```sh
mini-redis-cli.exe help
```

You can also specify host and port:
```sh
./mini-redis-cli --host 127.0.0.1 --port 6379 GET foo
```

> **⚠️ Important:**
> 
> The CLI **must use the same host and port as the server**. If you started the server with custom `--host` or `--port` values, use the same values when running the CLI. Otherwise, the connection will fail.


## Notes

- This project is **not recommended for production use**.
- The focus is on learning, not on performance or full Redis compatibility.
- The implementation covers only basic commands and simplified persistence.

## Credits

Developed as part of a personal study and inspired by the [Codecrafters Redis Go challenge](https://codecrafters.io/challenges/redis). 