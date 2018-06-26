package utils

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"os/exec"
	"strings"
)

// CommandSubstitution - Execute a command and return the output value. No exit on stdErr
func CommandSubstitution(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, _ := cmd.CombinedOutput()

	DebugOutput(stdoutStderr)

	return strings.TrimSpace(string(stdoutStderr))
}

// CommandSubSplitOutput - Execute a command and return the output value split into stdout and stderr. No exit on stdErr
func CommandSubSplitOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		DebugErrOutput(err)
		return "", err
	}

	DebugOutput(stdoutStderr)

	return strings.TrimSpace(string(stdoutStderr)), nil
}
