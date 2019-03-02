package system

import (
	"log"
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
		log.Fatal("Tokaido is currently only compatible with limited Linux distributions and MacOS High Sierra or newer.")
		return ""
	}
}
