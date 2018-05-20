package utils

import (
	"runtime"
)

// CheckOS - Check the users' operating system (runtime)
func CheckOS() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "osx"
	case "linux":
		return "linux"
	default:
		// TODO
		return "FatalError(Error(\"Tokaido is currently only compatible with Linux and Mac OSX operating systems\"))"
	}
}
