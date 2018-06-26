package utils

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

	"os/exec"
	"strings"
)

// BashStringCmd - Execute a bash command from a string `bash -c "(cmd)" with no log output`
func BashStringCmd(cmdStr string) string {
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, _ := cmd.CombinedOutput()

	DebugOutput(stdoutStderr)

	return strings.TrimSpace(string(stdoutStderr))
}

// BashStringSplitOutput - Execute a bash command splitting the resulting stdout and stderr
func BashStringSplitOutput(cmdStr string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		DebugErrOutput(err)
		return "", err
	}

	DebugOutput(stdoutStderr)

	return strings.TrimSpace(string(stdoutStderr)), nil
}
