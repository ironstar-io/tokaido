package utils

import (
	"github.com/ironstar-io/tokaido/conf"

	"github.com/fatih/color"
)

// DebugCmd ...
func DebugCmd(cmd string) {
	if conf.GetConfig().Tokaido.Debug == true && cmd != "" {
		color.Yellow(Timestamp()+": Running command: '%s'", cmd)
	}
}

// DebugString ...
func DebugString(output string) {
	if conf.GetConfig().Tokaido.Debug == true && len(output) > 0 {
		color.Yellow(Timestamp() + ": " + output)
	}
}

// DebugOutput ...
func DebugOutput(output []byte) {
	if conf.GetConfig().Tokaido.Debug == true && len(output) > 0 {
		color.Yellow("%s: %s", Timestamp(), output)
	}
}

// DebugErrOutput ...
func DebugErrOutput(output error) {
	if conf.GetConfig().Tokaido.Debug == true {
		color.Red("%s: %s", Timestamp(), output)
	}
}
