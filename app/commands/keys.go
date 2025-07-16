package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"mini-redis-go/app/server_config"
)

func handleKeys(args []string, _ map[string]Entry, config server_config.ServerConfig) string {
	if err := validateKeysCommand(args); err != nil {
		return RedisError(err.Error())
	}

	pattern := args[1]
	if pattern != "*" {
		return RedisEmptyArray
	}

	dbFilePath := filepath.Join(config.Dir, config.DBFileName)
	keys, err := extractRDBKeys(dbFilePath)
	if err != nil {
		return RedisEmptyArray
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

func extractRDBKeys(rdbFilePath string) ([]string, error) {
	rdbFile, err := os.Open(rdbFilePath)
	if err != nil {
		return nil, err
	}
	defer rdbFile.Close()

	var keys []string
	var opcodeBuffer [1]byte

	if _, err := io.ReadFull(rdbFile, make([]byte, RDB_HEADER_SIZE)); err != nil {
		return nil, err
	}
	if _, err := io.ReadFull(rdbFile, make([]byte, RDB_VERSION_SIZE)); err != nil {
		return nil, err
	}

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
				return nil, err
			}
			_, err = readString(rdbFile)
			if err != nil {
				return nil, err
			}
		} else if opcode == RDB_OPCODE_EXPIRETIME || opcode == RDB_OPCODE_EXPIRETIME_MS {
			_ = skipExpireTimeField(rdbFile, opcode)
			rdbFile.Read(opcodeBuffer[:])
			_, err := readString(rdbFile)
			if err != nil {
				return nil, err
			}
			err = skipValue(rdbFile, opcodeBuffer[0])
			if err != nil {
				return nil, err
			}
		} else if opcode == RDB_OPCODE_RESIZEDB {
			_, err := readLength(rdbFile)
			if err != nil {
				return nil, err
			}
			_, err = readLength(rdbFile)
			if err != nil {
				return nil, err
			}
		} else if opcode <= RDB_TYPE_HASH_ZIPLIST {
			keyName, err := readString(rdbFile)
			if err != nil {
				return nil, err
			}
			keys = append(keys, keyName)
			err = skipValue(rdbFile, opcode)
			if err != nil {
				return nil, err
			}
		}
	}
	return keys, nil
}

func readLength(reader io.Reader) (int, error) {
	var opcodeBuffer [1]byte
	if _, err := reader.Read(opcodeBuffer[:]); err != nil {
		return 0, err
	}

	opcode := opcodeBuffer[0]
	typ := (opcode & LENGTH_TYPE_MASK) >> 6
	switch typ {
	case LENGTH_6BIT:
		return read6BitLength(opcode), nil
	case LENGTH_14BIT:
		return read14BitLength(opcode, reader)
	case LENGTH_32BIT:
		return read32BitLength(reader)
	default:
		return 0, fmt.Errorf("unsupported length encoding")
	}
}

func read6BitLength(opcode byte) int {
	return int(opcode & LENGTH_6BIT_MASK)
}

func read14BitLength(opcode byte, reader io.Reader) (int, error) {
	var lowByteBuffer [1]byte
	if _, err := reader.Read(lowByteBuffer[:]); err != nil {
		return 0, err
	}
	return int(opcode&LENGTH_14BIT_MASK)<<8 | int(lowByteBuffer[0]), nil
}

func read32BitLength(reader io.Reader) (int, error) {
	var fourByteBuffer [4]byte
	if _, err := reader.Read(fourByteBuffer[:]); err != nil {
		return 0, err
	}
	return int(fourByteBuffer[0])<<24 | int(fourByteBuffer[1])<<16 | int(fourByteBuffer[2])<<8 | int(fourByteBuffer[3]), nil
}

func readString(reader io.Reader) (string, error) {
	var firstByte [1]byte
	if _, err := reader.Read(firstByte[:]); err != nil {
		return "", err
	}

	if firstByte[0] >= 0xC0 {
		return "", nil
	}

	length, err := readLengthFromFirstByte(reader, firstByte[0])
	if err != nil {
		return "", err
	}
	singleByteBuffer := make([]byte, length)
	if _, err := io.ReadFull(reader, singleByteBuffer); err != nil {
		return "", err
	}
	return string(singleByteBuffer), nil
}

func readLengthFromFirstByte(reader io.Reader, firstByte byte) (int, error) {
	opcode := firstByte
	typ := (opcode & LENGTH_TYPE_MASK) >> 6
	switch typ {
	case LENGTH_6BIT:
		return read6BitLength(opcode), nil
	case LENGTH_14BIT:
		return read14BitLength(opcode, reader)
	case LENGTH_32BIT:
		return read32BitLength(reader)
	default:
		return 0, fmt.Errorf("unsupported length encoding")
	}
}

func skipDatabaseIndexField(reader io.Reader) error {
	_, err := readLength(reader)
	return err
}

func skipExpireTimeField(reader io.Reader, opcode byte) error {
	var size int
	if opcode == RDB_OPCODE_EXPIRETIME {
		size = RDB_EXPIRETIME_SIZE
	} else {
		size = RDB_EXPIRETIME_MS_SIZE
	}
	_, err := io.ReadFull(reader, make([]byte, size))
	return err
}

func skipValue(reader io.Reader, valueType byte) error {
	switch valueType {
	case RDB_TYPE_STRING:
		_, err := readString(reader)
		return err
	default:
		return fmt.Errorf("unsupported value type: %d", valueType)
	}
}
