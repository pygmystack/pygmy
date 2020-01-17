package library

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/fubarhouse/pygmy-go/service/interface"
)

// NetworkCreate is part of a centralised abstraction of the Docker API
// and will create a Docker network with a specified configuration.
func NetworkCreate(network types.NetworkResource) error {
	return model.DockerNetworkCreate(&network)
}

// NetworkConnect is part of a centralised abstraction of the Docker API
// and will connect a created container to a docker network with a
// specified name.
func NetworkConnect(c Config, name string, containerName string) error {
	return model.DockerNetworkConnect(name, containerName)
}

// NetworkStatus will check the state of a Docker network to test if it has
// been created, and will return false if the network can not be found.
func NetworkStatus(network string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return false, err
	}
	networkResources, _ := cli.NetworkList(ctx, types.NetworkListOptions{})
	for _, Network := range networkResources {
		if Network.Name == network {
			return true, nil
		}
	}
	return false, errors.New(fmt.Sprintf("network %v not found\n", network))
}
