package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"mini-redis-go/app/server_config"
)

func handleGet(args []string, store map[string]Entry, config server_config.ServerConfig) string {
	if len(args) < 2 && args[0] == "GET" {
		return RedisError("wrong number of arguments for 'get' command")
	}

	key := args[1]

	entry, ok := store[key]
	if ok {
		if hasEntryExpired(entry) {
			delete(store, key)
			return "$-1\r\n"
		}
		return fmt.Sprintf("$%d\r\n%s\r\n", len(entry.Value), entry.Value)
	}

	dbFilePath := filepath.Join(config.Dir, config.DBFileName)
	value, found, expiry, err := GetRDBStringValueWithExpiry(dbFilePath, key)
	if err != nil {
		return RedisError("error reading RDB")
	}
	if !found {
		return "$-1\r\n"
	}
	if expiry > 0 && time.Now().UnixMilli() > expiry {
		return "$-1\r\n"
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
}

func GetRDBStringValueWithExpiry(rdbFilePath, searchKey string) (string, bool, int64, error) {
	rdbFile, err := os.Open(rdbFilePath)
	if err != nil {
		return "", false, 0, err
	}
	defer rdbFile.Close()

	if _, err := io.ReadFull(rdbFile, make([]byte, RDB_HEADER_SIZE)); err != nil {
		return "", false, 0, err
	}
	if _, err := io.ReadFull(rdbFile, make([]byte, RDB_VERSION_SIZE)); err != nil {
		return "", false, 0, err
	}

	var opcodeBuffer [1]byte
	var pendingExpiry int64
	for {
		if _, err := rdbFile.Read(opcodeBuffer[:]); err != nil {
			break
		}
		opcode := opcodeBuffer[0]

		if opcode == RDB_OPCODE_SELECTDB {
			_ = skipDatabaseIndexField(rdbFile)
		} else if opcode == RDB_OPCODE_EOF {
			break
		} else if opcode == RDB_OPCODE_AUX {
			_, err := readString(rdbFile)
			if err != nil {
				return "", false, 0, err
			}
			_, err = readString(rdbFile)
			if err != nil {
				return "", false, 0, err
			}
		} else if opcode == RDB_OPCODE_EXPIRETIME {
			var buf [4]byte
			if _, err := io.ReadFull(rdbFile, buf[:]); err != nil {
				return "", false, 0, err
			}
			sec := int64(buf[3])<<24 | int64(buf[2])<<16 | int64(buf[1])<<8 | int64(buf[0])
			pendingExpiry = sec * 1000
		} else if opcode == RDB_OPCODE_EXPIRETIME_MS {
			var buf [8]byte
			if _, err := io.ReadFull(rdbFile, buf[:]); err != nil {
				return "", false, 0, err
			}
			ms := int64(buf[7])<<56 | int64(buf[6])<<48 | int64(buf[5])<<40 | int64(buf[4])<<32 |
				int64(buf[3])<<24 | int64(buf[2])<<16 | int64(buf[1])<<8 | int64(buf[0])
			pendingExpiry = ms
		} else if opcode == RDB_OPCODE_RESIZEDB {
			_, err := readLength(rdbFile)
			if err != nil {
				return "", false, 0, err
			}
			_, err = readLength(rdbFile)
			if err != nil {
				return "", false, 0, err
			}
		} else if opcode <= RDB_TYPE_HASH_ZIPLIST {
			keyName, err := readString(rdbFile)
			if err != nil {
				return "", false, 0, err
			}
			value, err := readString(rdbFile)
			if err != nil {
				return "", false, 0, err
			}
			resultExpiry := pendingExpiry
			pendingExpiry = 0
			if keyName == searchKey {
				if resultExpiry > 0 && time.Now().UnixMilli() > resultExpiry {
					continue
				}
				return value, true, resultExpiry, nil
			}
		}
	}
	return "", false, 0, nil
}
