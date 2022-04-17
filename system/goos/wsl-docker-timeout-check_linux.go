package goos

import (
	"fmt"
	"os"
)

// CheckComposeTimeoutValue checks if COMPOSE_HTTP_TIMEOUT env var has been set or not.
// Errors occur when fallback to default value in WSL
func CheckComposeTimeoutValue() {
	if os.Getenv("COMPOSE_HTTP_TIMEOUT") == "" {
		fmt.Println(`The environment variable COMPOSE_HTTP_TIMEOUT has not been set on your system.

The value for this environment variable will generally default to 60 if it is unset, however due to a bug in WSL,
 this value isn't read correctly and with cause docker-compose commands to timeout.

To resolve this issue, please set COMPOSE_HTTP_TIMEOUT environment variable to a numeric value, we recommend 300.`)
	}
}
