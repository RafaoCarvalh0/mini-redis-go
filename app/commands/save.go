package commands

import (
	"mini-redis-go/app/server_config"
	"os"
	"path/filepath"
)

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
