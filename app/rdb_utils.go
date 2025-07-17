package main

import (
	"fmt"
	"io"
	"os"

	"mini-redis-go/app/commands"
	"mini-redis-go/app/server_config"
)

func loadDataFromRDB(config server_config.ServerConfig) (map[string]commands.Entry, error) {
	store := make(map[string]commands.Entry)

	if config.Dir != "" && config.DBFileName != "" {
		rdbFilePath := config.Dir + string(os.PathSeparator) + config.DBFileName

		loadedStore, err := loadRDBDataIntoStore(rdbFilePath, store)
		if err != nil {
			if os.IsNotExist(err) {
				createEmptyRDBFile(config, rdbFilePath)
				return store, nil
			}

			return store, fmt.Errorf("error loading rdb: %s", err.Error())
		}

		return loadedStore, nil
	}

	return store, fmt.Errorf("unexpected error while loading rdb")
}

func createEmptyRDBFile(config server_config.ServerConfig, rdbFilePath string) {
	if mkErr := os.MkdirAll(config.Dir, 0755); mkErr != nil {
		println("Error creating directory:", mkErr.Error())
		os.Exit(1)
	}

	f, createErr := os.Create(rdbFilePath)
	if createErr != nil {
		println("Error creating RDB file:", createErr.Error())
		os.Exit(1)
	}

	f.Close()
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

func loadRDBDataIntoStore(rdbFilePath string, store map[string]commands.Entry) (map[string]commands.Entry, error) {
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
			store[key] = commands.Entry{Value: value, ExpiryTime: 0}
		}
	}
	return store, nil
}
