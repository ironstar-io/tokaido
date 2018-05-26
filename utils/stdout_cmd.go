package utils

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"fmt"
	"log"
	"os/exec"
	"strings"
)

// StdoutCmd - Execute a command on the users' OS
func StdoutCmd(name string, args ...string) string {
	stdoutStderr := SilentStdoutCmd(name, args...)

	fmt.Printf("%s\n", stdoutStderr)

	return string(stdoutStderr)
}

// SilentStdoutCmd - Execute a command on the users' OS without logging result to the console
func SilentStdoutCmd(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Tokaido encountered a fatal error and had to stop at command '%s %s'\n%s", name, strings.Join(args, " "), stdoutStderr)
		log.Fatal(err)
	}

	return strings.TrimSpace(string(stdoutStderr))
}

// NoFatalStdoutCmd - Execute a command on the users' OS without exiting on stdoutError
func NoFatalStdoutCmd(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, _ := cmd.CombinedOutput()

	return strings.TrimSpace(string(stdoutStderr))
}
