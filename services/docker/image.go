package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// ImageExists - Check if an image exists
func ImageExists(image string) bool {
	dockerImage := utils.CommandSubstitution("docker", "images", "-q", image)

	if dockerImage == "" {
		return false
	}

	return true
}
