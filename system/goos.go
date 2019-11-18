package system

import (
	"log"
	"runtime"
)

// CheckOS - Check the users' operating system (runtime)
func CheckOS() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "macos"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		log.Fatal("Tokaido is currently only compatible with limited Linux distributions, Windows 10 Pro and MacOS High Sierra or newer.")
		return ""
	}
}
