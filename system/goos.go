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
	case "windows":
		return "windows"
	default:
		log.Fatal("Tokaido is currently only compatible with limited Linux distributions, Mac OSX and Windows operating systems")
		return ""
	}
}
