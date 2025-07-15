package server_config

import "os"

const DefaultRdbStorePath = "/tmp/redis_files"
const DefaultRdbFileName = "dump.rdb"

type ServerConfig struct {
	Dir        string
	DBFileName string
}

func GetConfigFromArgs() ServerConfig {
	var dir string
	var dbfilename string

	for i, arg := range os.Args {
		if arg == "--dir" && i+1 < len(os.Args) {
			dir = os.Args[i+1]
		}
		if arg == "--dbfilename" && i+1 < len(os.Args) {
			dbfilename = os.Args[i+1]
		}
	}
	return ServerConfig{
		Dir:        dir,
		DBFileName: dbfilename,
	}
}
