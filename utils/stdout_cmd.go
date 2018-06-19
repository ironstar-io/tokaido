package utils

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"

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

<<<<<<< HEAD
// StdoutStreamCmd - Execute a command on the users' OS and stream stdout
func StdoutStreamCmd(name string, args ...string) {
	DebugCmd(name + " " + strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
=======
// CommandSubSplitOutput - Execute a command and return the output value split into stdout and stderr. No exit on stdErr
func CommandSubSplitOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(stdoutStderr)), nil
}

// SilentBashStringCmd - Execute a bash command from a string `bash -c "(cmd)" with no log output`
func SilentBashStringCmd(cmdStr string) string {
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, _ := cmd.CombinedOutput()

	return strings.TrimSpace(string(stdoutStderr))
}

// BashStringSplitOutput - Execute a bash command splitting the resulting stdout and stderr
func BashStringSplitOutput(cmdStr string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = fs.WorkDir()
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(stdoutStderr)), nil
}

// BashStringCmd - Execute a bash command from a string `bash -c "(cmd)"`
func BashStringCmd(cmdStr string) string {
	stdoutStderr := SilentBashStringCmd(cmdStr)

	fmt.Printf("%s\n", stdoutStderr)

	return string(stdoutStderr)
>>>>>>> Begin build out for debug configuration
}
