package git

import (
	"log"

	"github.com/ironstar-io/tokaido/utils"
)

// AddAll - Run the `git add .` (stage all files) command
func AddAll() {
	_, err := utils.CommandSubSplitOutput("git", "add", ".")
	if err != nil {
		log.Fatal("Unable to run `git add .`:", err)
	}
}
