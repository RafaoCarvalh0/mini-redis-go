package protocol_parser

import (
	"fmt"
	"net"
	"strings"
)

func GetRESP2ArgsFromConn(conn net.Conn) ([]string, error) {
	var errResponse []string
	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		return errResponse, err
	}

	resp2Command := buffer[:n]
	args, err := parseRESP2Command(string(resp2Command))
	if err != nil {
		return errResponse, err
	}

	return args, nil
}

func parseRESP2Command(resp2Command string) ([]string, error) {
	lines := strings.Split(string(resp2Command), "\r\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("incomplete resp2 message")
	}
	if !strings.HasPrefix(lines[0], "*") {
		return nil, fmt.Errorf("not a resp2 array")
	}

	var result []string
	for i := 2; i < len(lines); i += 2 {
		if lines[i] == "" {
			break
		}
		result = append(result, lines[i])
	}

	return result, nil
}
