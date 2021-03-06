package ssh

import (
	"log"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/utils"
)

// ConnectCommand - Aliases `ssh <project-name> -C command`
func ConnectCommand(args []string) string {
	if len(args) == 0 {
		log.Fatal("At least one argument must be supplied to use this command")
	}

	cs := strings.Join(args, " ")
	pn := conf.GetConfig().Tokaido.Project.Name + ".tok"

	r, err := utils.CommandSubSplitOutput("ssh", []string{"-q", "-o UserKnownHostsFile=/dev/null", "-o StrictHostKeyChecking=no", pn, "-C", cs}...)
	if err != nil {
		log.Fatal(err)
	}

	utils.DebugString(r)

	return r
}

// ConnectCommandOutput - Aliases `ssh <project-name> -C command` and return any output stream as a value
func ConnectCommandOutput(args []string) string {
	if len(args) == 0 {
		log.Fatal("At least one argument must be supplied to use this command")
	}

	cs := strings.Join(args, " ")
	pn := conf.GetConfig().Tokaido.Project.Name + ".tok"

	r := utils.CommandSubstitution("ssh", []string{"-q", "-o UserKnownHostsFile=/dev/null", "-o StrictHostKeyChecking=no", pn, "-C", cs}...)

	utils.DebugString(r)

	return r
}

// StreamConnectCommand - Aliases `ssh <project-name> -C command` and streams std channels
func StreamConnectCommand(args []string) {
	if len(args) == 0 {
		log.Fatal("At least one argument must be supplied to use this command")
	}

	cs := strings.Join(args, " ")
	pn := conf.GetConfig().Tokaido.Project.Name + ".tok"

	utils.StreamOSCmd("ssh", []string{"-q", "-o UserKnownHostsFile=/dev/null", "-o StrictHostKeyChecking=no", pn, "-C", cs}...)
}
