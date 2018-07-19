package utils

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"

	"github.com/fatih/color"
)

// DebugCmd ...
func DebugCmd(cmd string) {
	if conf.GetConfig().Debug == true && cmd != "" {
		color.Yellow("Running command: '%s'", cmd)
	}
}

// DebugString ...
func DebugString(output string) {
	if conf.GetConfig().Debug == true && len(output) > 0 {
		color.Yellow("%s", output)
	}
}

// DebugOutput ...
func DebugOutput(output []byte) {
	if conf.GetConfig().Debug == true && len(output) > 0 {
		color.Yellow("%s", output)
	}
}

// DebugErrOutput ...
func DebugErrOutput(output error) {
	if conf.GetConfig().Debug == true {
		color.Red("%s", output)
	}
}
