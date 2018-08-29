package goos

import (
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	utils.CheckCmdHard("brew")

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
		_, perr := exec.LookPath("pip")
		if perr == nil {
			fmt.Println("    The dependency 'unison-fsmonitor' is missing. Tokaido will attempt to install it with pip")

			ut := filepath.Join(fs.WorkDir(), "unox.tar.gz")
			td := filepath.Join(fs.WorkDir(), "unox-0.2.0")
			_, cerr := utils.CommandSubSplitOutput("curl", "-o", ut, "https://codeload.github.com/hnsl/unox/tar.gz/0.2.0")
			if cerr != nil {
				fmt.Println(cerr)
				fsmonitorFatal()
			}
			_, gerr := utils.BashStringSplitOutput("gunzip -c " + ut + " | tar xopf -")
			if gerr != nil {
				fmt.Println(gerr)
				fsmonitorFatal()
			}
			_, uerr := utils.CommandSubSplitOutput("pip", "install", td)
			if gerr != nil {
				fmt.Println(uerr)
				fsmonitorFatal()
			}

			fs.Remove(ut)
			fs.RemoveAll(td)

			fmt.Println("  √  unison-fsmonitor")

			return
		}

		fsmonitorFatal()
	}

	// Temporarily disabled automated install through brew due to an issue with the install process
	// Outlined in this issue: https://github.com/ironstar-io/tokaido/issues/22

	// if err != nil {
	// 	fmt.Println("    unison-fsmonitor is missing. Tokaido will install it with Homebrew")
	// 	utils.StdoutCmd("brew", "tap", "eugenmayer/dockersync")
	// 	utils.StdoutCmd("brew", "install", "eugenmayer/dockersync/unox")
	// 	fmt.Println("  √  unison-fsmonitor")
	// }

	return
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
