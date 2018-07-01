package docker

import (
	"strings"

	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"bufio"
	"fmt"
	"os"
)

// ComposeStdout - Convenience method for docker-compose shell commands
func ComposeStdout(args ...string) {
	composeParams := composeArgs(args...)

	utils.StdoutCmd("docker-compose", composeParams...)
}

// ComposeResult - Convenience method for docker-compose shell commands returning a result
func ComposeResult(args ...string) string {
	composeParams := composeArgs(args...)

	return utils.CommandSubstitution("docker-compose", composeParams...)
}

func composeArgs(args ...string) []string {
	composeFile := []string{"-f", fs.WorkDir() + "/docker-compose.tok.yml"}

	return append(composeFile, args...)
}

// Up - Lift all containers in the compose file
func Up() {
	if ImageExists("tokaido/drush:latest") == false {
		fmt.Println(`🚡  First time lifting your containers? There's a few images to download, this might take some time.`)
	}

	ComposeStdout("up", "-d")
}

// Stop - Stop all containers in the compose file
func Stop() {
	ComposeStdout("stop")
}

// Down - Pull down all the containers in the compose file
func Down() {
	confirmDestroy := utils.ConfirmationPrompt(`
🔥  This will also destroy the database inside your Tokaido environment. Are you sure?`, "n")

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

// Ps - Print the container status to the console
func Ps() {
	ComposeStdout("ps")
}

// StatusCheck ...
func StatusCheck() {
	rawStatus := ComposeResult("ps")

	unavailableContainers := false
	scanner := bufio.NewScanner(strings.NewReader(rawStatus))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Name") || strings.Contains(scanner.Text(), "------") {
			continue
		} else if !strings.Contains(scanner.Text(), "Up") {
			unavailableContainers = true
		}
	}

	if unavailableContainers == true {
		fmt.Println(`
😓 Tokaido containers are not working properly

It appears that some or all of the Tokaido containers are offline.

You can try to fix this by running 'tok up', or by running 'tok repair'.
	`)
		os.Exit(1)
	}

	fmt.Println(`
✅  All containers are running`)
}
