package library

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/fubarhouse/pygmy-go/service/interface"
)

func NetworkCreate(c Config, name string) error {
	return model.DockerNetworkCreate(name, c.Networks[name].Config)
}

func NetworkConnect(c Config, name string, containerName string) error {
	return model.DockerNetworkConnect(name, containerName)
}

func NetworkStatus(network string) (bool, error) {
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
