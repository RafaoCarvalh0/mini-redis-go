package main

import (
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/protocol_parser"
	"github.com/codecrafters-io/redis-starter-go/app/server_config"
)

func main() {
	store := make(map[string]commands.Entry)
	config := server_config.GetConfigFromArgs()

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		os.Exit(1)
	}

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
