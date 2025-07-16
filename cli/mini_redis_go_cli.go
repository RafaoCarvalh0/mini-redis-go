package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	hostEnv := os.Getenv("MINI_REDIS_HOST")
	portEnv := os.Getenv("MINI_REDIS_PORT")

	hostDefault := hostEnv
	if hostDefault == "" {
		hostDefault = "127.0.0.1"
	}
	portDefault := portEnv
	if portDefault == "" {
		portDefault = "6379"
	}

	host := flag.String("host", hostDefault, "Redis server host")
	port := flag.String("port", portDefault, "Redis server port")
	flag.Parse()

	if flag.NArg() == 0 || flag.Arg(0) == "help" {
		printHelp()
		os.Exit(0)
	}

	address := net.JoinHostPort(*host, *port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Could not connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	argv := flag.Args()
	command := buildRESP2Command(argv)
	_, err = conn.Write([]byte(command))
	if err != nil {
		fmt.Println("Error sending command:", err)
		os.Exit(1)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading response (timeout or connection closed):", err)
		os.Exit(1)
	}
	fmt.Print(string(buf[:n]))
}

func printHelp() {
	fmt.Println("mini-redis-cli - Available commands:")
	fmt.Println("  PING                   - Test the connection with the server")
	fmt.Println("  ECHO <message>         - Echo back the provided message")
	fmt.Println("  SET <key> <value> [PX <milliseconds>] - Set a value for a key, optionally with expiration in ms (PX)")
	fmt.Println("  GET <key>              - Get the value of a key")
	fmt.Println("  CONFIG <subcommand>    - Manage server configuration")
	fmt.Println("  KEYS <pattern>         - List keys matching the pattern")
	fmt.Println("  SAVE                   - Save the current dataset to disk")
	fmt.Println("  help                   - Show this help message")
}

func buildRESP2Command(argv []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", len(argv)))
	for _, arg := range argv {
		sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}
	return sb.String()
}
