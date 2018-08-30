package docker

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/utils"

	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
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
	cs := utils.StdoutCmd("docker-compose", "-f", filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/docker-compose.tok.yml"), "port", containerName, containerPort)

	p := strings.Split(cs, ":")
	if len(p) <= 1 {
		log.Fatal("The required container '" + containerName + "' is not running. Check the status of your containers with `tok ps` and retry `tok up`.  Exiting...")
		return ""
	}

	return p[1]
}

// GetContainerName ...
func GetContainerName(name string) (string, error) {
	cs, err := utils.BashStringSplitOutput("docker-compose -f " + filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/docker-compose.tok.yml") + " ps " + name + " | grep " + name)
	if err != nil {
		return "", err
	}

	cn := strings.Split(cs, " ")[0]
	if cn == "" {
		return "", errors.New("Unable to find the container for " + name + ". Is it currently running?")
	}

	return cn, nil
}

// GetContainerIP ...
func GetContainerIP(name string) (string, error) {
	cn, err := GetContainerName(name)
	if err != nil {
		return "", err
	}

	ia, err := utils.BashStringSplitOutput(`docker inspect ` + cn + ` | grep "IPAddress\": \"1"`)
	if err != nil {
		return "", err
	}

	for _, s := range strings.Split(ia, `"`) {
		if strings.Contains(s, "1") {
			return s, nil
		}
	}

	return "", errors.New("Unable to find the internal IP for container " + name + ". Is it currently running?")
}

// KillContainer - Kill an individual container
func KillContainer(container string) {
	utils.StdoutCmd("docker-compose", "-f", filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/docker-compose.tok.yml"), "kill", container)
}

// UpContainer - Lift an individual container
func UpContainer(container string) {
	utils.StdoutCmd("docker-compose", "-f", filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/docker-compose.tok.yml"), "up", "-d", container)
}
