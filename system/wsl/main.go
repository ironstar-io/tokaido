package wsl

import (
	"io/ioutil"
	"strings"

	"github.com/ironstar-io/tokaido/system/fs"
)

// IsWSL - Return true if OS is using Windows Subsystem for Linux (WSL)
func IsWSL() bool {
	if fs.CheckExists("/proc/version") == false {
		return false
	}

	b, err := ioutil.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	s := string(b)

	return strings.Contains(s, "Microsoft")
}
