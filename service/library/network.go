package library

import (
	"context"
	"fmt"
	networktypes "github.com/docker/docker/api/types/network"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pygmystack/pygmy/service/interface/docker"
)

// NetworkCreate is part of a centralised abstraction of the Docker API
// and will create a Docker network with a specified configuration.
func NetworkCreate(network networktypes.Inspect) error {
	return docker.DockerNetworkCreate(&network)
}

// NetworkConnect is part of a centralised abstraction of the Docker API
// and will connect a created container to a docker network with a
// specified name.
func NetworkConnect(network string, containerName string) error {
	return docker.DockerNetworkConnect(network, containerName)
}

// NetworkStatus will check the state of a Docker network to test if it has
// been created, and will return false if the network can not be found.
func NetworkStatus(network string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	cli.NegotiateAPIVersion(ctx)
	if err != nil {
		return false, err
	}
	networkResources, _ := cli.NetworkList(ctx, types.NetworkListOptions{})
	for _, Network := range networkResources {
		if Network.Name == network {
			return true, nil
		}
	}
	return false, fmt.Errorf("network %v not found\n", network)
}
