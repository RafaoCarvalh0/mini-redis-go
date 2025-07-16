package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"mini-redis-go/app/commands"
	"mini-redis-go/app/protocol_parser"
	"mini-redis-go/app/server_config"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dir := flag.String("dir", "../", "RDB directory (default: repo root)")
	dbfilename := flag.String("dbfilename", "example_dump.rdb", "RDB file name (default: template_dump.rdb)")
	hostEnv := os.Getenv("MINI_REDIS_HOST")
	portEnv := os.Getenv("MINI_REDIS_PORT")

	hostDefault := hostEnv
	if hostDefault == "" {
		hostDefault = "0.0.0.0"
	}
	portDefault := portEnv
	if portDefault == "" {
		portDefault = "6379"
	}

	host := flag.String("host", hostDefault, "Server host")
	port := flag.String("port", portDefault, "Server port")
	flag.Parse()

	config := server_config.GetServerConfig(dir, dbfilename)

	store, err := loadRDB(config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("File", *dbfilename, "loaded successfully")

	address := net.JoinHostPort(*host, *port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Server listening on:", address)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn, store, config)
	}
}

func handleConnection(conn net.Conn, store map[string]commands.Entry, config server_config.ServerConfig) {
	defer conn.Close()
	for {
		args, err := protocol_parser.GetRESP2ArgsFromConn(conn)
		if err != nil {
			break
		}

		if resp, ok := commands.HandleCommand(args, store, config); ok {
			conn.Write([]byte(resp))
		}
	}
}
