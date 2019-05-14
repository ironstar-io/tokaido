package git

import (
	"log"

	"github.com/ironstar-io/tokaido/utils"
)

// Init - Run the `git init` command
func Init(dir string) {
	_, err := utils.CommandSubSplitOutputContext(dir, "git", "init")
	if err != nil {
		log.Fatal("Unable to run `git init`:", err)
	}
}
