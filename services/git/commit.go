package git

import (
	"log"

	"github.com/ironstar-io/tokaido/utils"
)

// Commit - Run the `git commit -m <message>` command
func Commit(message string) {
	_, err := utils.CommandSubSplitOutput("git", "commit", "-m", message)
	if err != nil {
		log.Fatal("Unable to run `git commit -m <message>`:", err)
	}
}
