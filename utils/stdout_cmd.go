package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
)

// StdoutCmd - Execute a command on the users' OS
func StdoutCmd(name string, args ...string) string {
	DebugCmd(name + " " + strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Dir = conf.GetProjectPath()
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

// StdoutStreamCmdDebug - Execute a command on the users' OS and stream stdout for debug only
func StdoutStreamCmdDebug(name string, args ...string) {
	CommandSubSplitOutputContext(conf.GetProjectPath(), name, args...)
}

// StdoutStreamCmdDebugContext - Execute a command on the users' OS and stream stdout for debug only
func StdoutStreamCmdDebugContext(directory, name string, args ...string) {
	DebugCmd(name + " " + strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Dir = directory

	if conf.GetConfig().Tokaido.Debug == true {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.Run()
}
