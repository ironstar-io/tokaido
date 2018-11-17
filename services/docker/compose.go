package docker

import (
	"bufio"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/utils"
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
	composeFile := []string{"-f", filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/docker-compose.tok.yml")}

	return append(composeFile, args...)
}

// CreateVolume - Create a new docker volume with a name
func CreateVolume(name string) {
	o := utils.BashStringCmd(`docker volume ls | grep ` + name)
	if o == name {
		utils.DebugString("Existing volume found with the name '" + name + "', not creating a new one")
		return
	}

	utils.StdoutCmd("docker", "volume", "create", name)
}

// CreateDatabaseVolume will create a database volume if it doesn't already exist
func CreateDatabaseVolume() {
	n := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_mysql_database"
	CreateVolume(n)
}

// CreateComposerCacheVolume will create a composer cache volume if it doesn't already exist
func CreateComposerCacheVolume() {
	n := "tok_composer_cache"
	CreateVolume(n)
}

// Up - Lift all containers in the compose file
func Up() {
	ComposeStdout("up", "-d")
}

// UpMulti - Lift a specific list of containers
func UpMulti(args []string) {
	c := append([]string{"up", "-d"}, args...)

	ComposeStdout(c...)
}

// Stop - Stop all containers in the compose file
func Stop() {
	fmt.Println()
	w := console.SpinStart("Tokaido is stopping your containers!")

	ComposeStdout("stop")

	console.SpinPersist(w, "🚉", "Tokaido stopped your containers successfully!")
}

// Down - Pull down all the containers in the compose file
func Down() {
	w := console.SpinStart("Tokaido is taking down your containers!")

	ComposeStdout("down")

	console.SpinPersist(w, "🚉", "Successfully destroyed your Tokaido environment")
}

// PrintLogs - Print all logs or the container logs to the console
func PrintLogs(args []string) {
	ls := append([]string{"logs"}, args...)

	fmt.Println(ComposeResult(ls...))
}

// Ps - Print the container status to the console
func Ps() {
	fmt.Println(ComposeResult("ps"))
}

// Exec - Execute a command inside the service container
func Exec(args []string) {
	cm := append([]string{"exec", "-T"}, args...)
	fmt.Println(ComposeResult(cm...))
}

// StatusCheck ...
func StatusCheck() error {
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

		return errors.New("Tokaido containers are not running correctly")
	}

	return nil
}
