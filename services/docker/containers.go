package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/utils"
	"golang.org/x/net/context"
)

// getContainerState returns the state of a specified container
func getContainerState(name string) (state string, err error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.36"))
	if err != nil {
		return "", err
	}

	filter := filters.NewArgs()
	filter.Add("name", conf.GetConfig().Tokaido.Project.Name+"_"+name)

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filter,
	})
	if err != nil {
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
