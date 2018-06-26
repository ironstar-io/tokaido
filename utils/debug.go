package utils

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"

	"fmt"
)

// DebugOutput ...
func DebugOutput(output []byte) {
	if conf.GetConfig().Debug == true {
		fmt.Printf("Debug: %s\n", output)
	}
}

// DebugErrOutput ...
func DebugErrOutput(output error) {
	if conf.GetConfig().Debug == true {
		fmt.Printf("Debug: %s\n", output)
	}
}
