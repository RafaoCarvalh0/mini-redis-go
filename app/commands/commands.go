package commands

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"mini-redis-go/app/server_config"
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

func hasEntryExpired(entry Entry) bool {
	return entry.ExpiryTime > 0 && time.Now().UnixMilli() > entry.ExpiryTime
}

func LoadRDBDataIntoStore(rdbFilePath string, store map[string]Entry) (map[string]Entry, error) {
	rdbFile, err := os.Open(rdbFilePath)
	if err != nil {
		return store, err
	}
	defer rdbFile.Close()

	if _, err := io.ReadFull(rdbFile, make([]byte, RDB_HEADER_SIZE)); err != nil {
		return store, err
	}
	if _, err := io.ReadFull(rdbFile, make([]byte, RDB_VERSION_SIZE)); err != nil {
		return store, err
	}

	var opcodeBuffer [1]byte
	for {
		if _, err := rdbFile.Read(opcodeBuffer[:]); err != nil {
			break
		}
		opcode := opcodeBuffer[0]

		if opcode == RDB_OPCODE_EOF {
			break
		}
		if opcode == RDB_TYPE_STRING {
			key, err := readString(rdbFile)
			if err != nil {
				return store, err
			}
			value, err := readString(rdbFile)
			if err != nil {
				return store, err
			}
			store[key] = Entry{Value: value, ExpiryTime: 0}
		}
	}
	return store, nil
}
