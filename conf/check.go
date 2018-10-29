package conf

import (
	"log"

	"github.com/ironstar-io/tokaido/constants"
)

// ValidProjectRoot - Check if the project root was found, if not log and exit
func ValidProjectRoot() {
	if GetConfig().Tokaido.Project.Name == constants.ProjectRootNotFound {
		log.Fatal(constants.ProjectNotFoundError)
	}
}
