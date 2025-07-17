package commands

import (
	"fmt"
	"strings"

	"mini-redis-go/app/server_config"
)

type Entry struct {
	Value      string
	ExpiryTime int64
}

var commandHandlers = map[string]func([]string, *map[string]Entry, server_config.ServerConfig) string{
	"PING":   handlePing,
	"ECHO":   handleEcho,
	"SET":    handleSet,
	"GET":    handleGet,
	"CONFIG": handleConfig,
	"KEYS":   handleKeys,
	"SAVE":   handleSave,
	"DEL":    handleDel,
}

func HandleCommand(args []string, store *map[string]Entry, config server_config.ServerConfig) (string, bool) {
	command := strings.ToUpper(args[0])

	if handler, isCommandHandled := commandHandlers[command]; isCommandHandled {
		return handler(args, store, config), true
	}

	return RedisError(fmt.Sprintf("unknown command '%s'", args[0])), false
}
