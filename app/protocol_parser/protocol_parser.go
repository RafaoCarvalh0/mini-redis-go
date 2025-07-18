package protocol_parser

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func GetRESP2ArgsFromConn(conn net.Conn) ([]string, error) {
	emptySlice := []string{}
	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		return emptySlice, err
	}

	resp2Command := buffer[:n]
	args, err := parseRESP2Command(string(resp2Command))
	if err != nil {
		return emptySlice, err
	}

	return args, nil
}

func parseRESP2Command(resp2Command string) ([]string, error) {
	var parsedCommand []string

	lines := strings.Split(string(resp2Command), "\r\n")

	if len(lines) == 0 || !strings.HasPrefix(lines[0], "*") || !isLengthCorrect(lines) {
		return parsedCommand, fmt.Errorf("invalid resp2 command")
	}

	for i := 2; i < len(lines); i += 2 {
		if lines[i] == "" {
			break
		}
		parsedCommand = append(parsedCommand, lines[i])
	}

	return parsedCommand, nil
}

func isLengthCorrect(lines []string) bool {
	fmt.Println(lines)

	var expectedStrLength int

	for i := 1; i < len(lines); i++ {
		currentCharacter := lines[i]

		if currentCharacter == "" {
			return true
		}

		afterCut, foundPrefix := strings.CutPrefix(currentCharacter, "$")

		if foundPrefix {
			lengthInt, err := strconv.Atoi(afterCut)
			if err != nil {
				return false
			}

			expectedStrLength = lengthInt
			continue
		}

		if len(currentCharacter) != expectedStrLength {
			return false
		}
	}

	return false
}
