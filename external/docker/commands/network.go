package commands

import (
	"context"
	"fmt"

	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/networks"
)

// NetworkCreate is part of a centralised abstraction of the Docker API
// and will create a Docker network with a specified configuration.
// TODO: To be migrated to the Docker runtime provider.
func NetworkCreate(ctx context.Context, cli *client.Client, network networktypes.Inspect) error {
	return networks.Create(ctx, cli, &network)
}

// NetworkConnect is part of a centralised abstraction of the Docker API
// and will connect a created container to a docker network with a
// specified name.
// TODO: To be migrated to the Docker runtime provider.
func NetworkConnect(ctx context.Context, cli *client.Client, network string, containerName string) error {
	return networks.Connect(ctx, cli, network, containerName)
}

// NetworkStatus will check the state of a Docker network to test if it has
// been created, and will return false if the network can not be found.
// TODO: To be migrated to the Docker runtime provider.
func NetworkStatus(network string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	cli.NegotiateAPIVersion(ctx)
	if err != nil {
		return false, err
	}
	networkResources, _ := cli.NetworkList(ctx, networktypes.ListOptions{})
	for _, Network := range networkResources {
		if Network.Name == network {
			return true, nil
		}
	}
	return false, fmt.Errorf("network %v not found\n", network)
}
