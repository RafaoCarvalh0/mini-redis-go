package commands

import (
	"fmt"
	"mini-redis-go/app/server_config"
	"strings"
)

func handleEcho(args []string, _ *map[string]Entry, _ server_config.ServerConfig) string {
	if strings.ToUpper(args[0]) == "ECHO" {
		if len(args) == 2 {
			msg := args[1]
			resp := fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)
			return resp
		}

		return RedisError("wrong number of arguments for 'echo' command")
	}

	return ""
}
