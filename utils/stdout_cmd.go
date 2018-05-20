package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// StdoutCmd - Execute a command on the users' OS
func StdoutCmd(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Tokaido encountered a fatal error and had to stop at command '%s %s'\n%s", name, strings.Join(args, " "), stdoutStderr)
		return FatalError(err)
	}

	fmt.Printf("%s\n", stdoutStderr)
	return string(stdoutStderr)
}
