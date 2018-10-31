package docker

import (
	"github.com/ironstar-io/tokaido/utils"
)

// PullImages - Pull all images in compose file
func PullImages() {
	ComposeStdout("pull", "-q")
}

// ImageExists - Check if an image exists
func ImageExists(image string) bool {
	dockerImage := utils.CommandSubstitution("docker", "images", "-q", image)

	if dockerImage == "" {
		return false
	}

	return true
}
