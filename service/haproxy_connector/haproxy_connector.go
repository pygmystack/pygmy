package haproxy_connector

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/fubarhouse/pygmy-go/service/interface"
)

// Connect will connect a specified container to a specified network.
func Connect(containerName string, network string) error {
	if s, _ := Connected(containerName, network); !s {
		return model.DockerNetworkConnect(network, containerName)
	}
	return nil
}

// Connected will check if a container is already connected to a network
// with a given name provided as input.
func Connected(containerName string, network string) (bool, error) {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return false, err
	}
	x, err := cli.NetworkInspect(ctx, network)
	if err != nil {
		return false, err
	}
	for _, container := range x.Containers {
		if container.Name == containerName {
			return true, nil
		}
	}
	return false, nil
}
