package utils

import (
	"os/exec"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
)

// CommandSubstitution - Execute a command and return the output value. No exit on stdErr
func CommandSubstitution(name string, args ...string) string {
	DebugCmd(name + " " + strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Dir = conf.GetProjectPath()
	stdoutStderr, _ := cmd.CombinedOutput()

	DebugOutput(stdoutStderr)

	return strings.TrimSpace(string(stdoutStderr))
}

// CommandSubstitutionExitCode - Execute a command and return 0 if successful, or 1 if failed
// We can't return the true exit code on MacOS as this isn't a reliable method and will have if
// `name` is not available in the users' path (see: https://stackoverflow.com/questions/10385551/get-exit-code-go)
func CommandSubstitutionExitCode(name string, args ...string) (exitCode int) {
	DebugCmd(name + " " + strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Dir = conf.GetProjectPath()
	err := cmd.Run()

	e, ok := err.(*exec.ExitError)
	if ok && !e.Success() {
		return 0
	}

	return 1
}

// CommandSubSplitOutput - Execute a command and return the output value split into stdout and stderr. No exit on stdErr
func CommandSubSplitOutput(name string, args ...string) (string, error) {
	return CommandSubSplitOutputContext(conf.GetProjectPath(), name, args...)
}

// CommandSubSplitOutputContext - Execute a command and return the output value split into stdout and stderr from a directory context. No exit on stdErr
func CommandSubSplitOutputContext(directory, name string, args ...string) (string, error) {
	DebugCmd(name + " " + strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Dir = directory
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		DebugErrOutput(err)
		return "", err
	}

	DebugOutput(stdoutStderr)

	return strings.TrimSpace(string(stdoutStderr)), nil
}
