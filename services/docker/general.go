package docker

import (
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/spf13/viper"

	"github.com/ironstar-io/tokaido/utils"

	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var defaultContainers = map[string]string{
	"drush":   "22",
	"fpm":     "9000",
	"haproxy": "8443",
	"mysql":   "3306",
	"nginx":   "8082",
	"unison":  "5000",
}

var optionalContainers = map[string]string{
	"solr":    "8983",
	"redis":   "6379",
	"adminer": "8080",
	"mailhog": "8025",
}

// PrintPorts - Print a port map of all containers, or a single container
func PrintPorts(containers []string) {
	if len(containers) == 0 {
		data := make([][]string, len(defaultContainers))

		for k, v := range defaultContainers {
			data = append(data, []string{k, LocalPort(k, v)})
		}

		for k, v := range optionalContainers {
			if viper.GetString("services."+k+".enabled") == "true" {
				data = append(data, []string{k, LocalPort(k, v)})
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Container", "Current Port"})
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.AppendBulk(data)
		table.Render()

		return
	}

	target := containers[0]
	if defaultContainers[target] == "" {
		fmt.Println("There aren't any containers matching the name '" + target + "'")
		return
	}

	fmt.Println(LocalPort(target, defaultContainers[target]))
}

// LocalPort - Return the local port of a container
func LocalPort(containerName string, containerPort string) string {
	// Example return: "unison:32757"
	containerStr := utils.StdoutCmd("docker-compose", "-f", filepath.Join(fs.WorkDir(), "/docker-compose.tok.yml"), "port", containerName, containerPort)

	return strings.Split(containerStr, ":")[1]
}

// KillContainer - Kill an individual container
func KillContainer(container string) {
	utils.StdoutCmd("docker-compose", "-f", filepath.Join(fs.WorkDir(), "/docker-compose.tok.yml"), "kill", container)
}

// UpContainer - Lift an individual container
func UpContainer(container string) {
	utils.StdoutCmd("docker-compose", "-f", filepath.Join(fs.WorkDir(), "/docker-compose.tok.yml"), "up", "-d", container)
}
