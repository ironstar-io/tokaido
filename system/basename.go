package system

import (
	"path/filepath"
)

// Basename - Return the name of the directory the binary is being run in
func Basename() string {
	return filepath.Base(WorkDir())
}
