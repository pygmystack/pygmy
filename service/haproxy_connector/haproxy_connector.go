package haproxy_connector

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/fubarhouse/pygmy/v1/service/interface"
)

func Connect(containerName string, network string) error {
	if s, _ := Connected(containerName, network); !s {
		return model.DockerNetworkConnect(network, containerName)
	}
	return nil
}

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
