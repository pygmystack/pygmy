package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	model "github.com/fubarhouse/pygmy/service/interface"
)

func Create() error {
	return model.DockerNetworkCreate("amazeeio-network")
}

func Connect() error {
	return model.DockerNetworkConnect("amazeeio-network", "amazeeio-haproxy")
}

func Status() (bool, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return false, err
	}
	networkResources, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	for _, Network := range networkResources {
		if Network.Name == "amazeeio-network" {
			return true, nil
		}
	}
	return false, errors.New(fmt.Sprintf("network amazeeio-network not found\n"))
}
