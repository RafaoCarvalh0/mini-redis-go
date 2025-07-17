package commands

import (
	"fmt"
	"mini-redis-go/app/server_config"
	"strings"
	"time"
)

func handleSet(args []string, store *map[string]Entry, _ server_config.ServerConfig) string {
	if len(args) > 0 && strings.ToUpper(args[0]) == "SET" {
		fmt.Println(args)
		if len(args) != 3 {
			return RedisError("wrong number of arguments for 'set' command")
		}
		key := args[1]
		value := args[2]
		expiry := int64(0)

		if len(args) > 4 && strings.ToUpper(args[3]) == "PX" {
			ms, err := parseInt64(args[4])
			if err != nil {
				return RedisError("PX value is not an integer")
			}
			expiry = time.Now().UnixMilli() + ms
		}

		(*store)[key] = Entry{Value: value, ExpiryTime: expiry}
		return "+OK\r\n"
	}

	return ""
}
