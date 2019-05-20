package docker

import (
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

// ComposeExitCode - Convenience method for docker-compose shell commands returning just the bash exit code
func ComposeExitCode(args ...string) int {
	composeParams := composeArgs(args...)

	return utils.CommandSubstitutionExitCode("docker-compose", composeParams...)
}

func composeArgs(args ...string) []string {
	composeFile := []string{"-f", filepath.Join(conf.GetProjectPath(), "/docker-compose.tok.yml")}

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

// CreateSiteVolume will create a Site volume if it doesn't already exist
func CreateSiteVolume() {
	n := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_tokaido_site"
	CreateVolume(n)
}

// CreateComposerCacheVolume will create a composer cache volume if it doesn't already exist
func CreateComposerCacheVolume() {
	n := "tok_composer_cache"
	CreateVolume(n)
}

// Up - Lift all containers in the compose file
func Up() {
	ComposeStdout("up", "--remove-orphans", "-d")
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

// ExecWithExitCode - Execute a command inside the admin container and return an error if it failed to run
func ExecWithExitCode(args []string) (err error) {
	cm := append([]string{"exec", "-T"}, args...)
	result := ComposeExitCode(cm...)

	if result != 0 {
		return errors.New("Command '%s' returned an unknown error response")
	}

	return nil

}

// StatusCheck tests for containers 'running' state.
// If one container is specified and true if that container is 'running'
// If no container is specified, all containers must be 'running' or will return false
func StatusCheck(container, project string) (ok bool) {
	cl := GetContainerList()

	// If the user specified a single container to check, only go for that one
	if len(container) > 1 {
		state, err := getContainerState(container, project)
		if err != nil {
			panic(err)
		}

		if state == "running" {
			return true
		}
		return false
	}

	// The user didn't specify a container so we'll check them all
	utils.DebugString("Checking container state for [" + strings.Join(cl, ",") + "]")
	for _, v := range cl {
		state, err := getContainerState(v, project)
		if state != "running" {
			if err != nil {
				utils.DebugErrOutput(err)
			}
			return false
		}
	}

	return true
}
