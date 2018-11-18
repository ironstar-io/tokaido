package goos

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	CheckBrew()
	CheckAndBrewInstall("unison")
	CheckDockersync()
}

// CheckBrew ...
func CheckBrew() *string {
	_, err := exec.LookPath("brew")
	if err != nil {
		fmt.Println("    Homebrew isn't installed. Tokaido will install it")
		utils.StdoutCmd("/usr/bin/ruby", "-e", "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\"")
		fmt.Println("  √  brew")
	}

	return nil
}

// CheckDockersync - Root executable
func CheckDockersync() {
	_, err := exec.LookPath("unison-fsmonitor")
	if err != nil {
		fmt.Println("    The dependency 'unison-fsmonitor' is missing. Tokaido will install it with Homebrew")
		err = stdUnisonInstall()
		if err != nil {
			fmt.Println(err)
			err = altUnisonInstall()
			if err != nil {
				fmt.Println(err)
				fsmonitorFatal()
			}
		}

		fmt.Println("  √  unison-fsmonitor")
	}

	return
}

func stdUnisonInstall() error {
	_, err := utils.CommandSubSplitOutput("brew", "tap", "eugenmayer/dockersync")
	if err != nil {
		return err
	}

	_, err = utils.CommandSubSplitOutput("brew", "install", "eugenmayer/dockersync/unox")
	if err != nil {
		return err
	}

	return nil
}

func altUnisonInstall() error {
	_, err := exec.LookPath("pip")
	if err == nil {
		fmt.Println("    The dependency 'unison-fsmonitor' was unable to be installed with brew. Tokaido will attempt to install it with pip")

		ut := filepath.Join(fs.WorkDir(), "unox.tar.gz")
		td := filepath.Join(fs.WorkDir(), "unox-0.2.0")
		_, err = utils.CommandSubSplitOutput("curl", "-o", ut, "https://codeload.github.com/hnsl/unox/tar.gz/0.2.0")
		if err != nil {
			return err
		}
		_, err = utils.BashStringSplitOutput("gunzip -c " + ut + " | tar xopf -")
		if err != nil {
			return err
		}
		_, err = utils.CommandSubSplitOutput("pip", "install", td)
		if err != nil {
			return err
		}

		fs.Remove(ut)
		fs.RemoveAll(td)

		fmt.Println("  √  unison-fsmonitor")

		return nil
	}

	return err
}

func fsmonitorFatal() {
	log.Fatal(`
Tokaido is unable to install the required dependency 'unison-fsmonitor' for you automatically.

There is an issue outlining this raised at: https://github.com/ironstar-io/tokaido/issues/22

It is still possible to run Tokaido, but you will need to install the dependency manually.

This requires git, and the Python package management tool, pip, to be installed.
$ easy_install pip

Then you should be able to install 'unison-fsmonitor'
$ git clone https://github.com/hnsl/unox
$ cd unox
$ pip install .

This workaround is temporary, sorry for the inconvenience.

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
