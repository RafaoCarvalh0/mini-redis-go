package commands

import (
	"fmt"
	"mini-redis-go/app/server_config"
	"strings"
)

func handleGet(args []string, store map[string]Entry, config server_config.ServerConfig) string {
	if len(args) < 2 && strings.ToUpper(args[0]) == "GET" {
		return RedisError("wrong number of arguments for 'get' command")
	}

	key := args[1]

	entry, ok := store[key]
	if ok {
		if hasEntryExpired(entry) {
			delete(store, key)
			return "$-1\r\n"
		}
		return fmt.Sprintf("$%d\r\n%s\r\n", len(entry.Value), entry.Value)
	}

	return "$-1\r\n"
}
