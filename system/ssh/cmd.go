package ssh

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/utils"

	"log"
	"strings"
)

// ConnectCommand - Aliases `ssh <project-name> -C command`
func ConnectCommand(args []string) string {
	if len(args) == 0 {
		log.Fatal("At least one argument must be supplied to use this command")
	}

	cs := strings.Join(args, " ")
	pn := conf.GetConfig().Tokaido.Project.Name + ".tok"

	return utils.CommandSubstitution("ssh", []string{"-q", "-o UserKnownHostsFile=/dev/null", "-o StrictHostKeyChecking no", pn, "-C", cs}...)
}
