//go:build !windows
// +build !windows

package goos

import (
	"golang.org/x/sys/unix"
)

// Writable - Check if a file/folder is writable
func Writable(path string) bool {
	return unix.Access(path, unix.W_OK) == nil
}
