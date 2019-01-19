package goos

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/ironstar-io/tokaido/utils"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	CheckBrew()
	CheckPip()
	CheckAndBrewInstall("unison")
	CheckAndBrewInstall("unison-fsmonitor")
	CheckDockersync()
}

// CheckBrew ...
func CheckBrew() *string {
	_, err := exec.LookPath("brew")
	if err != nil {
		log.Fatal("    Homebrew is missing. You can get it from https://brew.sh")
	}
	fmt.Println("  √  brew")

	return nil
}

// CheckPip ...
func CheckPip() {
	_, err := exec.LookPath("pip")
	if err != nil {
		log.Fatal("    Python Pip isn't installed. You can install with 'sudo easy_install pip'")
	}
	fmt.Println("  √  pip")
}

// CheckDockersync - Root executable
func CheckDockersync() {
	_, err := exec.LookPath("unison-fsmonitor")
	if err != nil {
		fmt.Println("    The dependency 'unison-fsmonitor' is missing. Tokaido will install it with Homebrew")
		err = unisonInstall()
		if err != nil {
			fsmonitorFatal()
		}

		fmt.Println("  √  unison-fsmonitor")
	}

	return
}

func unisonInstall() error {
	_, err := utils.CommandSubSplitOutput("brew", "tap", "ironstar-io/tap")
	if err != nil {
		return err
	}

	_, err = utils.CommandSubSplitOutput("brew", "install", "ironstar-io/tap/unison-fsmonitor")
	if err != nil {
		return err
	}

	_, err = utils.CommandSubSplitOutput("pip", "install", "watchdog")
	if err != nil {
		return err
	}

	return nil
}

func fsmonitorFatal() {
	log.Fatal(`
Tokaido is unable to install the required dependency 'unison-fsmonitor' for you automatically. 

This dependency is required for Tokaido's background sync functionality. 

You can manually download the unison-fsmonitor Python script with the following commands:

curl https://github.com/ironstar-io/unison-fsmonitor/releases/download/0.0.1/unison-fsmonitor.py -o unison-fsmonitor 
chmod +x unison-fsmonitor
sudo cp unison-fsmonitor /usr/local/bin/unison-fsmonitor
pip install watchdog

Exiting...
	`)
}

// CheckAndBrewInstall - Root executable
func CheckAndBrewInstall(program string) *string {
	_, err := exec.LookPath(program)
	if err != nil {
		fmt.Println("     " + program + " isn't installed. Tokaido will install it with Homebrew")
		utils.StdoutCmd("brew", "install", program)
		fmt.Println("  √ ", program)
	}

	return nil
}
