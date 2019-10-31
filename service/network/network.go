package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/fubarhouse/pygmy/v0/service/interface"
)

func Create(network string) error {
	return model.DockerNetworkCreate(network)
}

func Connect(containerName string, network string) error {
	return model.DockerNetworkConnect(network, containerName)
}

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
