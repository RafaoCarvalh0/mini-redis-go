package server_config

const DefaultRdbStorePath = "/tmp/redis_files"
const DefaultRdbFileName = "dump.rdb"

type ServerConfig struct {
	Dir        string
	DBFileName string
}

func GetServerConfig(dir, dbfilename *string) ServerConfig {

	return ServerConfig{
		Dir:        *dir,
		DBFileName: *dbfilename,
	}
}
