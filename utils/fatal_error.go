package utils

import (
	"log"
)

// FatalError - Log a fatal error and os.Exit(1)
var FatalError = func(err error) string {
	if err != nil {
		log.Fatal(err)
	}

	return ""
}
