package commands

import (
	"fmt"
	"strings"
)

func RedisError(errMsg string) string {
	finalMsg := errMsg

	if strings.Trim(errMsg, " ") == "" {
		finalMsg = "unknown error"
	}

	return fmt.Sprintf("-ERR %s\r\n", finalMsg)
}
