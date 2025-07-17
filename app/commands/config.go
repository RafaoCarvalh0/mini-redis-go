package commands

import (
	"fmt"
	"mini-redis-go/app/server_config"
	"strings"
)

func handleConfig(args []string, _ *map[string]Entry, config server_config.ServerConfig) string {
	if err := validateConfigCommand(args); err != nil {
		return RedisError(err.Error())
	}

	switch strings.ToUpper(args[2]) {
	case "DIR":
		storePath := config.Dir
		resp := fmt.Sprintf("*2\r\n$3\r\ndir\r\n$%d\r\n%s\r\n", len(storePath), storePath)
		return resp
	case "DBFILENAME":
		rdbFilename := config.DBFileName
		resp := fmt.Sprintf("*2\r\n$9\r\ndbfilename\r\n$%d\r\n%s\r\n", len(rdbFilename), rdbFilename)
		return resp
	}

	return RedisEmptyArray
}

func validateConfigCommand(args []string) error {
	if len(args) != 3 || strings.ToUpper(args[1]) != "GET" {
		return fmt.Errorf("wrong number of arguments for 'config' command")
	}

	subcmd := strings.ToUpper(args[2])
	switch subcmd {
	case "DIR", "DBFILENAME":
		return nil
	default:
		return fmt.Errorf("unknown command '%s'", args[2])
	}
}
