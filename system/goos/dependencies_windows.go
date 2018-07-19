package goos

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
	"os/exec"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	checkChoco()
	checkUnison()
}

// CheckChoco ...
func checkChoco() {
	_, err := exec.LookPath("choco")
	if err != nil {
		fmt.Println("\nChocolately isn't installed. Tokaido will install it")
		utils.PowershellCmd("Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))")
		fmt.Println("Chocolatey installed successfully")
	}
}

// checkUnison ...
func checkUnison() *string {
	_, err := exec.LookPath("unison")
	if err != nil {
		fmt.Println("Unison isn't installed. Tokaido will install it with Chocolately")
		utils.StdoutCmd("choco", "install", "unison")
		fmt.Println("Unison installed successfully")
	}

	return nil
}

// CheckAndChocoInstall ...
func CheckAndChocoInstall(program string) *string {
	_, err := exec.LookPath(program)
	if err != nil {
		fmt.Println("     " + program + " isn't installed. Tokaido will install it with Homebrew")
		utils.StdoutCmd("choco", "install", program)
		fmt.Println(program + " installed successfully")
	}

	return nil
}
