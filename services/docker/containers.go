package docker

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/ironstar-io/tokaido/utils"
	"golang.org/x/net/context"
)

// getContainerState returns the state of a specified container
func getContainerState(name, project string) (state string, err error) {
	dcli := GetAPIClient()
	project = strings.ToLower(project)

	filter := filters.NewArgs()
	cn := fmt.Sprintf("%s-%s", project, name)
	filter.Add("name", cn)

	utils.DebugString(fmt.Sprintf("looking for container name: %s", cn))

	containers, err := dcli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filter,
	})
	if err != nil {
		fmt.Println("Docker API failed: ")
		fmt.Println(err)
		return "", err
	}

	if len(containers) > 1 {
		return "", fmt.Errorf("error looking up container state for container %s. Received %d matching containers when only wanting 1", name, len(containers))
	}

	if len(containers) < 1 {
		utils.DebugErrOutput(fmt.Errorf("error looking up container state for container %s. Could not find a container by that name", name))
		return "", nil
	}

	if containers[0].State == "" {
		return "", fmt.Errorf("error looking up container state for container %s. Container state was empty", name)
	}

	return containers[0].State, nil
}

// GetContainers returns all containers for a project
func GetContainers(project string) *[]types.Container {
	dcli := GetAPIClient()
	project = strings.ToLower(project)

	filter := filters.NewArgs()
	filter.Add("name", project+"_")

	containers, err := dcli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filter,
	})
	if err != nil {
		panic(err)
	}

	return &containers
}

// GetContainer returns a docker Container object for the specified container, if it exists
func GetContainer(name, project string) types.Container {
	dcli := GetAPIClient()
	project = strings.ToLower(project)

	filter := filters.NewArgs()
	cn := project + "-" + name
	filter.Add("name", cn)

	containers, err := dcli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filter,
	})
	if err != nil {
		utils.DebugString(fmt.Sprintf("Warning: received error retrieving container ["+cn+"]: %v", err))
		return types.Container{}
	}

	if len(containers) != 1 {
		return types.Container{}
	}

	return containers[0]
}
