package commands

import "fmt"

func RedisError(errMsg string) string {
	return fmt.Sprintf("-ERR %s\r\n", errMsg)
}
