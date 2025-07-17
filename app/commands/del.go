package commands

import (
	"fmt"
	"mini-redis-go/app/server_config"
	"strings"
)

func handleDel(args []string, store *map[string]Entry, config server_config.ServerConfig) string {
	if len(args) < 2 && strings.ToUpper(args[0]) == "DEL" {
		return RedisError("wrong number of arguments for 'del' command")
	}

	deletedKeysCount := 0
	for i := 1; i < len(args); i++ {
		currentKey := args[i]

		if _, ok := (*store)[currentKey]; ok {
			delete(*store, currentKey)
			deletedKeysCount++
		}
	}

	return fmt.Sprintf(":%d\r\n", deletedKeysCount)
}
