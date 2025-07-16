package main

import (
	"fmt"
	"os"

	"mini-redis-go/app/commands"
	"mini-redis-go/app/server_config"
)

func loadRDB(config server_config.ServerConfig) (map[string]commands.Entry, error) {
	store := make(map[string]commands.Entry)
	fmt.Println(config)

	if config.Dir != "" && config.DBFileName != "" {
		rdbFilePath := config.Dir + string(os.PathSeparator) + config.DBFileName

		loadedStore, err := commands.LoadRDBDataIntoStore(rdbFilePath, store)
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
