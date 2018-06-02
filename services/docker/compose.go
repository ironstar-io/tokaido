package docker

import (
	"fmt"
)

// Up - Lift all containers in the compose file
func Up() {
	fmt.Println(`
🚡  First time lifting your containers? There's a few images to download, this might take some time.
	`)

	ComposeStdout("up", "-d")
}

// Stop - Stop all containers in the compose file
func Stop() {
	ComposeStdout("stop")
}

// Down - Pull down all the containers in the compose file
func Down() {
	ComposeStdout("down")
}

// Status - Print the container status to the console
func Status() {
	ComposeStdout("ps")
}
