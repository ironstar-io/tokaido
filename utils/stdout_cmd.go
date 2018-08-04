package utils

import (
	"github.com/ironstar-io/tokaido/system/fs"

	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// StdoutCmd - Execute a command on the users' OS
func StdoutCmd(name string, args ...string) string {
	DebugCmd(name + " " + strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Tokaido encountered a fatal error and had to stop at command '%s %s'\n%s", name, strings.Join(args, " "), stdoutStderr)
		log.Fatal(err)
	}

	DebugOutput(stdoutStderr)

	return strings.TrimSpace(string(stdoutStderr))
}

// StdoutStreamCmd - Execute a command on the users' OS and stream stdout
func StdoutStreamCmd(name string, args ...string) {
	DebugCmd(name + " " + strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
