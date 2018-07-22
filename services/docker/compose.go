package docker

import (
	"strings"

	"bitbucket.org/ironstar/tokaido-cli/system/console"
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
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
	composeFile := []string{"-f", filepath.Join(fs.WorkDir(), "/docker-compose.tok.yml")}

	return append(composeFile, args...)
}

// Up - Lift all containers in the compose file
func Up() {
	ComposeStdout("up", "-d")
}

// Stop - Stop all containers in the compose file
func Stop() {
	fmt.Println()
	w := wow.New(os.Stdout, spin.Get(spin.Dots), `   Tokaido is stopping your containers!`)
	w.Start()

	ComposeStdout("stop")

	w.PersistWith(spin.Spinner{Frames: []string{"🚉"}}, `  Tokaido stopped your containers successfully!`)
}

// Down - Pull down all the containers in the compose file
func Down() {
	confirmDestroy := utils.ConfirmationPrompt(`🔥  This will also destroy the database inside your Tokaido environment. Are you sure?`, "n")

	if confirmDestroy == false {
		console.Println(`🍵  Exiting without change`, "")
		return
	}

	fmt.Println()
	w := wow.New(os.Stdout, spin.Get(spin.Dots), `   Tokaido is pulling down your containers!`)
	w.Start()

	ComposeStdout("down")

	w.PersistWith(spin.Spinner{Frames: []string{"🚉"}}, `  Tokaido destroyed containers successfully!`)
}

// Ps - Print the container status to the console
func Ps() {
	fmt.Println(ComposeResult("ps"))
}

// StatusCheck ...
func StatusCheck() {
	rawStatus := ComposeResult("ps")

	unavailableContainers := false
	foundContainers := false
	scanner := bufio.NewScanner(strings.NewReader(rawStatus))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Name") || strings.Contains(scanner.Text(), "------") || strings.Contains(scanner.Text(), "cannot find the path specified") {
			continue
		} else if !strings.Contains(scanner.Text(), "Up") {
			unavailableContainers = true
			foundContainers = true
		}
		foundContainers = true
	}

	if unavailableContainers == true || foundContainers == false {
		console.Println(`
😓 Tokaido containers are not working properly`, "")
		fmt.Println(`
It appears that some or all of the Tokaido containers are offline.

View the status of your containers with 'tok ps'

You can try to fix this by running 'tok up', or by running 'tok repair'.
	`)
		os.Exit(1)
	}

	console.Println(`
✅  All containers are running`, "√")
}
