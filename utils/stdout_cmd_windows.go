// +build windows

package utils

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PowershellCmd - Execute a command in powershell
func PowershellCmd(args ...string) string {
	psc := filepath.Join(os.Getenv("SYSTEMROOT"), "System32/WindowsPowerShell/v1.0/powershell.exe")
	DebugCmd(psc + " " + strings.Join(args, " "))

	cmd := exec.Command(psc, args...)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Tokaido encountered a fatal error and had to stop at command '%s %s'\n%s", psc, strings.Join(args, " "), stdoutStderr)
		log.Fatal(err)
	}

	DebugOutput(stdoutStderr)

	return strings.TrimSpace(string(stdoutStderr))
}
