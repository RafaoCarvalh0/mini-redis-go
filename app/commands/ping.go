package commands

import "mini-redis-go/app/server_config"

func handlePing(args []string, _ *map[string]Entry, _ server_config.ServerConfig) string {
	if len(args) > 0 && args[0] == "PING" {
		return "+PONG\r\n"
	}
	return ""
}
