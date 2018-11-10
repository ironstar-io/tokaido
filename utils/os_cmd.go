package utils

import (
	"os"
	"os/exec"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
)

// StreamOSCmd - Execute a command outputting OS stdout/stdin/stderr directly to the console
func StreamOSCmd(name string, args ...string) {
	DebugCmd(name + " " + strings.Join(args, " "))

	c := exec.Command(name, args...)
	c.Dir = conf.GetConfig().Tokaido.Project.Path
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()
}
