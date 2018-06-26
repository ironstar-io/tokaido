package docker

import (
	"bitbucket.org/ironstar/tokaido-cli/system/fs"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var containerPorts = map[string]string{
	"drush":   "22",
	"fpm":     "9000",
	"haproxy": "8443",
	"mysql":   "3306",
	"nginx":   "8082",
	"solr":    "8983",
	"unison":  "5000",
}

// PrintPorts - Print a port map of all containers, or a single container
func PrintPorts(containers []string) {
	if len(containers) == 0 {
		data := make([][]string, len(containerPorts))
		for k, v := range containerPorts {
			data = append(data, []string{k, LocalPort(k, v)})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Container", "Current Port"})
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.AppendBulk(data)
		table.Render()

		return
	}

	target := containers[0]
	if containerPorts[target] == "" {
		fmt.Println("There aren't any containers matching the name '" + target + "'")
		return
	}

	fmt.Println(LocalPort(target, containerPorts[target]))
}

// LocalPort - Return the local port of a container
func LocalPort(containerName string, containerPort string) string {
	// Example return: "unison:32757"
	containerStr := utils.StdoutCmd("docker-compose", "-f", fs.WorkDir()+"/docker-compose.tok.yml", "port", containerName, containerPort)

	return strings.Split(containerStr, ":")[1]
}

// KillContainer - Kill an individual container
func KillContainer(container string) {
	utils.StdoutCmd("docker-compose", "-f", fs.WorkDir()+"/docker-compose.tok.yml", "kill", container)
}

// UpContainer - Lift an individual container
func UpContainer(container string) {
	utils.StdoutCmd("docker-compose", "-f", fs.WorkDir()+"/docker-compose.tok.yml", "up", "-d", container)
}
