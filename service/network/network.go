// Network provides an abstraction of the Docker API to create, connect and check docker network state.
// It does not disconnect, remove or check connection between the network and individual containers.
package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/fubarhouse/pygmy-go/service/model"
)

// Create is an abstraction for Docker network creation.
func Create(network string) error {
	return model.DockerNetworkCreate(network)
}

// Connect is an abstraction for connecting a container to a Docker network.
func Connect(containerName string, network string) error {
	return model.DockerNetworkConnect(network, containerName)
}

// Create is an abstraction for checking if a container is connected
// to a Docker network.
func Status(network string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return false, err
	}
	networkResources, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	for _, Network := range networkResources {
		if Network.Name == network {
			return true, nil
		}
	}
	return false, errors.New(fmt.Sprintf("network %v not found\n", network))
}
