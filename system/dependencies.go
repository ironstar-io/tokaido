package system

import (
	"fmt"
	"os"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/goos"

	aurora "github.com/logrusorgru/aurora"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	if conf.GetConfig().Tokaido.Dependencychecks == true {
		goos.CheckDependencies()
	}

	term := os.Getenv("TERM")
	if term != "screen-256color" && term != "xterm-256color" {
		fmt.Println(aurora.Yellow("Warning: The terminal '" + term + "' is not supported. Please switch to screen-256color or xterm-256color"))
	}
}
