package utils

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"

	"fmt"
)

// DebugCmd ...
func DebugCmd(cmd string) {
	if conf.GetConfig().Debug == true && cmd != "" {
		fmt.Printf("\033[33mRunning command: '%s'\033[0m\n", cmd)
	}
}

// DebugString ...
func DebugString(output string) {
	if conf.GetConfig().Debug == true && len(output) > 0 {
		fmt.Printf("\033[33m%s\033[0m\n", output)
	}
}

// DebugOutput ...
func DebugOutput(output []byte) {
	if conf.GetConfig().Debug == true && len(output) > 0 {
		fmt.Printf("\033[33m%s\033[0m\n", output)
	}
}

// DebugErrOutput ...
func DebugErrOutput(output error) {
	if conf.GetConfig().Debug == true {
		fmt.Printf("\033[31m%s\033[0m\n", output)
	}
}
