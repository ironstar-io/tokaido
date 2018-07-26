package utils

import (
	"time"
)

// Timestamp - Return current timestamp in RFC3339 format
func Timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
