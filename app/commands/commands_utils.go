package commands

import (
	"fmt"
	"time"
)

func parseInt64(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func hasEntryExpired(entry Entry) bool {
	return entry.ExpiryTime > 0 && time.Now().UnixMilli() > entry.ExpiryTime
}
