package docker

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/wsl"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

var defaultContainers = map[string]string{
	"drush":   "22",
	"fpm":     "9000",
	"haproxy": "8443",
	"mysql":   "3306",
	"nginx":   "8082",
	"varnish": "8081",
}

var optionalContainers = map[string]string{
	"solr":    "8983",
	"redis":   "6379",
	"adminer": "8080",
	"mailhog": "8025",
}

// PrintPorts - Print a port map of all containers, or a single container
func PrintPorts(containers []string) {
	// Return all containers if no specific container was provided
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

	// Return a single container
	target := containers[0]
	if len(defaultContainers[target]) > 0 {
		fmt.Println(LocalPort(target, defaultContainers[target]))
	} else if len(optionalContainers[target]) > 0 {
		fmt.Println(LocalPort(target, optionalContainers[target]))
	} else {
		fmt.Println("There aren't any containers matching the name '" + target + "'")
	}

	return
}

// LocalPort - Return the local port of a container
func LocalPort(containerName string, containerPort string) string {
	// Example return: "unison:32757"
	cs := utils.StdoutCmd("docker-compose", "-f", filepath.Join(conf.GetProjectPath(), "docker-compose.tok.yml"), "port", containerName, containerPort)

	p := strings.Split(cs, ":")
	if len(p) <= 1 {
		log.Fatal("The required container '" + containerName + "' is not running. Check the status of your containers with `tok ps` and retry `tok up`.  Exiting...")
		return ""
	}

	return p[1]
}

// GetContainerName ...
func GetContainerName(name string) (string, error) {
	cs, err := utils.BashStringSplitOutput("docker-compose -f " + filepath.Join(conf.GetProjectPath(), "docker-compose.tok.yml") + " ps " + name + " | grep " + name)
	if err != nil {
		return "", err
	}

	cn := strings.Split(cs, " ")[0]
	if cn == "" {
		return "", errors.New("Unable to find the container for " + name + ". Is it currently running?")
	}

	return cn, nil
}

// GetNetworkName takes a project name and returns a matching network name
// that complies with the same sanitisation that docker compose uses when
// it auto-generates networks from compose files
func GetNetworkName(pn string) string {
	return strings.ToLower(strings.Replace(pn, ".", "", -1)) + "_default"
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

// GetContainerIPFromProject returns the container IP for a given containerName in a given projectName
func GetContainerIPFromProject(containerName, projectName string) (string, error) {
	fullContainerName := fmt.Sprintf("%s_%s_1", projectName, containerName)

	ia, err := utils.BashStringSplitOutput(`docker inspect ` + fullContainerName + ` | grep "IPAddress\": \"1"`)
	if err != nil {
		return "", err
	}

	for _, s := range strings.Split(ia, `"`) {
		if strings.Contains(s, "1") {
			return s, nil
		}
	}

	return "", errors.New("Unable to find the internal IP for container " + containerName + ". Is it currently running?")

}

// GetContainerList returns a list of all configured containers in the current project context
func GetContainerList() []string {
	// List of containers to be checked when no single container is specifieds
	cl := []string{"fpm", "nginx", "varnish", "haproxy", "syslog", "drush", "mysql"}

	// Add additional optional containers to the checklist if they are in use
	if conf.GetConfig().Services.Kishu.Enabled {
		cl = append(cl, "kishu")
	}

	if conf.GetConfig().Services.Memcache.Enabled {
		cl = append(cl, "memcache")
	}

	if conf.GetConfig().Services.Redis.Enabled {
		cl = append(cl, "redis")
	}

	if conf.GetConfig().Services.Solr.Enabled {
		cl = append(cl, "solr")
	}

	if conf.GetConfig().Global.Syncservice == "fusion" {
		cl = append(cl, "sync")
	}

	return cl

}

// KillContainer - Kill an individual container
func KillContainer(container string) {
	utils.StdoutCmd("docker-compose", "-f", filepath.Join(conf.GetProjectPath(), "docker-compose.tok.yml"), "kill", container)
}

// UpContainer - Lift an individual container
func UpContainer(container string) {
	utils.StdoutCmd("docker-compose", "-f", filepath.Join(conf.GetProjectPath(), "docker-compose.tok.yml"), "up", "-d", container)
}

// DeleteVolume - Delete an external docker volume
func DeleteVolume(name string) {
	utils.StdoutCmd("docker", "volume", "rm", name)
}

// GetDockerHost - Return the value of the docker host that should be used for this system
func GetDockerHost() string {
	if wsl.IsWSL() {
		return "tcp://localhost:2375"
	}

	return client.DefaultDockerHost
}

// GetAPIClient - Returns a Docker API Client
func GetAPIClient() (dcli *client.Client) {
	h := GetDockerHost()

	dcli, err := client.NewClientWithOpts(client.WithVersion("1.36"), client.WithHost(h))
	if err != nil {
		log.Fatalf("Unable to obtain a Docker API client. Is Docker running?: %v", dcli)
	}

	return
}
