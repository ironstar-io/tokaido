package goos

import (
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	utils.CheckCmdHard("choco")
	checkUnison()
}

// checkUnison ...
func checkUnison() {
	_, err := exec.LookPath("unison")
	if err != nil {
		fmt.Println("Unison isn't installed. Tokaido will install it with Chocolately")

		utils.StdoutCmd("choco", "install", "unison", "-y")

		// Unison doesn't install as unison.exe in the chocolatey bin, this copies the installed binary so it's accessible from your PATH
		var ue string
		cbin := filepath.Join(os.Getenv("SYSTEMDRIVE")+"/", "ProgramData", "chocolatey", "bin")
		err := filepath.Walk(cbin, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if strings.Contains(path, "unison") && strings.Contains(path, "text") {
				ue = path
				return io.EOF
			}
			return nil
		})

		if err == io.EOF {
			fs.Copy(ue, filepath.Join(cbin, "unison.exe"))
			fmt.Println("Unison installed successfully")
		} else {
			fmt.Println("Unison may not have installed on your system correctly, this may cause issues with Tokaido initialization")
		}
	}
}

// CheckAndChocoInstall ...
func CheckAndChocoInstall(program string) *string {
	_, err := exec.LookPath(program)
	if err != nil {
		fmt.Println("     " + program + " isn't installed. Tokaido will install it with Homebrew")
		utils.StdoutCmd("choco", "install", program, "-y")
		fmt.Println(program + " installed successfully")
	}

	return nil
}
