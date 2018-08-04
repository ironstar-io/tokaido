package git

import (
	"fmt"
	"os"

	"github.com/ironstar-io/tokaido/system/fs"
)

// CheckGitRepo checks if we're running with a Git repo, fails if not
func CheckGitRepo() {
	if !fs.CheckExists(fs.WorkDir() + "/.git") {
		fmt.Println("🤷‍  Tokaido needs to be run in a Git repository root directory but could not find a .git directory here.")
		fmt.Println()
		os.Exit(1)
	}
}
