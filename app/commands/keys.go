package commands

import (
	"fmt"
	"mini-redis-go/app/server_config"
)

func handleKeys(args []string, store map[string]Entry, config server_config.ServerConfig) string {
	if err := validateKeysCommand(args); err != nil {
		return RedisError(err.Error())
	}

	pattern := args[1]
	if pattern != "*" {
		return RedisEmptyArray
	}

	var keys []string
	for key := range store {
		keys = append(keys, key)
	}

	resp := fmt.Sprintf("*%d\r\n", len(keys))
	for _, keyName := range keys {
		resp += fmt.Sprintf("$%d\r\n%s\r\n", len(keyName), keyName)
	}
	return resp
}

func validateKeysCommand(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("wrong number of arguments for 'keys' command")
	}

	if args[1] != "*" {
		return fmt.Errorf("unknown command %s", args[1])
	}

	return nil
}
