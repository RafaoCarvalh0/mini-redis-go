package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/server_config"
)

type Entry struct {
	Value      string
	ExpiryTime int64
}

var commandHandlers = map[string]func([]string, map[string]Entry, server_config.ServerConfig) string{
	"PING":   handlePing,
	"ECHO":   handleEcho,
	"SET":    handleSet,
	"GET":    handleGet,
	"CONFIG": handleConfig,
	"KEYS":   handleKeys,
	"SAVE":   handleSave,
}

func HandleCommand(args []string, store map[string]Entry, config server_config.ServerConfig) (string, bool) {
	if handler, isCommandHandled := commandHandlers[strings.ToUpper(args[0])]; isCommandHandled {
		return handler(args, store, config), true
	}

	return RedisError(fmt.Sprintf("unknown command '%s'", args[0])), false
}

func handlePing(args []string, _ map[string]Entry, _ server_config.ServerConfig) string {
	if len(args) > 0 && args[0] == "PING" {
		return "+PONG\r\n"
	}
	return ""
}

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

func handleSet(args []string, store map[string]Entry, _ server_config.ServerConfig) string {
	if len(args) > 0 && args[0] == "SET" {
		if len(args) < 3 {
			return RedisError("wrong number of arguments for 'set' command\r\n")
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
		store[key] = Entry{Value: value, ExpiryTime: expiry}
		return "+OK\r\n"
	}
	return ""
}

func hasEntryExpired(entry Entry) bool {
	return entry.ExpiryTime > 0 && time.Now().UnixMilli() > entry.ExpiryTime
}

func handleConfig(args []string, _ map[string]Entry, config server_config.ServerConfig) string {
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

	return "*0\r\n"
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

func handleSave(_ []string, _ map[string]Entry, config server_config.ServerConfig) string {
	if config.Dir == "" || config.DBFileName == "" {
		return RedisError("no directory or dbfilename provided")
	}
	if err := os.MkdirAll(config.Dir, 0755); err != nil {
		return RedisError("could not create directory: " + err.Error())
	}
	dbFilePath := filepath.Join(config.Dir, config.DBFileName)
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		f, err := os.Create(dbFilePath)
		if err != nil {
			return RedisError("could not create snapshot: " + err.Error())
		}
		defer f.Close()
	}
	return "+OK\r\n"
}
