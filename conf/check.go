package conf

import (
	"log"

	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
)

// ValidProjectRoot - Check if the project root was found, if not log and exit
func ValidProjectRoot() {
	if fs.ProjectRoot() == constants.ProjectRootNotFound {
		log.Fatal(constants.ProjectNotFoundError)
	}
}
