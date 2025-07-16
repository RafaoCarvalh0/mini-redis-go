package commands

import (
	"encoding/binary"
	"mini-redis-go/app/server_config"
	"os"
	"path/filepath"
)

const (
	rdbHeader      = "REDIS"
	rdbVersion     = "0001"
	length6BitMax  = 0x40
	length14BitMax = 0x4000
	length32BitTag = 0x80
)

func handleSave(_ []string, store map[string]Entry, config server_config.ServerConfig) string {
	if config.Dir == "" || config.DBFileName == "" {
		return RedisError("no directory or dbfilename provided")
	}

	if err := os.MkdirAll(config.Dir, 0755); err != nil {
		return RedisError("could not create directory: " + err.Error())
	}

	dbFilePath := filepath.Join(config.Dir, config.DBFileName)
	if err := writeRDBSnapshot(store, dbFilePath); err != nil {
		return RedisError("could not save snapshot: " + err.Error())
	}

	return "+OK\r\n"
}

func writeRDBSnapshot(store map[string]Entry, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write([]byte(rdbHeader)); err != nil {
		return err
	}
	if _, err := f.Write([]byte(rdbVersion)); err != nil {
		return err
	}

	for key, entry := range store {
		if entry.ExpiryTime > 0 {
			if _, err := f.Write([]byte{RDB_OPCODE_EXPIRETIME_MS}); err != nil {
				return err
			}
			if err := binary.Write(f, binary.LittleEndian, entry.ExpiryTime); err != nil {
				return err
			}
		}
		if _, err := f.Write([]byte{RDB_TYPE_STRING}); err != nil {
			return err
		}
		if err := writeRDBString(f, key); err != nil {
			return err
		}
		if err := writeRDBString(f, entry.Value); err != nil {
			return err
		}
	}

	if _, err := f.Write([]byte{RDB_OPCODE_EOF}); err != nil {
		return err
	}
	return nil
}

func writeRDBString(f *os.File, s string) error {
	length := int64(len(s))
	if err := writeRDBLength(f, length); err != nil {
		return err
	}
	if _, err := f.Write([]byte(s)); err != nil {
		return err
	}
	return nil
}

func writeRDBLength(f *os.File, length int64) error {
	if length < length6BitMax {
		b := byte(length)
		if _, err := f.Write([]byte{b}); err != nil {
			return err
		}
		return nil
	}
	if length < length14BitMax {
		b1 := encodeRDB14BitLengthPrefix(length)
		b2 := encodeRDB14BitLengthSuffix(length)
		if _, err := f.Write([]byte{b1, b2}); err != nil {
			return err
		}
		return nil
	}

	b := []byte{length32BitTag}
	if _, err := f.Write(b); err != nil {
		return err
	}

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(length))
	if _, err := f.Write(buf); err != nil {
		return err
	}
	return nil
}

func encodeRDB14BitLengthPrefix(length int64) byte {
	return byte(((length >> 8) & LENGTH_14BIT_MASK) | length6BitMax)
}

func encodeRDB14BitLengthSuffix(length int64) byte {
	return byte(length & RDB_OPCODE_EOF)
}
