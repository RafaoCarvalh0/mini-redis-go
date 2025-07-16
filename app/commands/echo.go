package commands

import (
	"fmt"
	"mini-redis-go/app/server_config"
)

func handleEcho(args []string, _ map[string]Entry, _ server_config.ServerConfig) string {
	if len(args) > 0 && args[0] == "ECHO" {
		if len(args) > 1 {
			msg := args[1]
			resp := fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)
			return resp
		}
		return RedisError("wrong number of arguments for 'echo' command")
	}
	return ""
}
