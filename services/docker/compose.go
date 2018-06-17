package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
)

// Up - Lift all containers in the compose file
func Up() {
	fmt.Println(`🚡  First time lifting your containers? There's a few images to download, this might take some time.`)

	ComposeStdout("up", "-d")
}

// Stop - Stop all containers in the compose file
func Stop() {
	ComposeStdout("stop")
}

// Down - Pull down all the containers in the compose file
func Down() {
	confirmDestroy := utils.ConfirmationPrompt(`
Are you sure? Destroying your containers will delete the database inside Tokaido`, "n")

	if confirmDestroy == false {
		fmt.Println(`
🍵  Exiting without change
		`)
		return
	}

	fmt.Println(`
🚅  Tokaido is pulling down your containers!
	`)

	ComposeStdout("down")

	fmt.Println(`
🚉  Tokaido destroyed containers successfully!
	`)
}

// Status - Print the container status to the console
func Status() {
	ComposeStdout("ps")
}
