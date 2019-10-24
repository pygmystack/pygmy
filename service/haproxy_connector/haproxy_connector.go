package haproxy_connector

import (
	"context"

	"github.com/docker/docker/client"
	model "github.com/fubarhouse/pygmy/service/interface"
)

func Connect() error {
	if s, _ := Connected(); !s {
		return model.DockerNetworkConnect("amazeeio-network", "amazeeio-haproxy")
	}
	return nil
}

func Connected() (bool, error) {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return false, err
	}
	x, err := cli.NetworkInspect(ctx, "amazeeio-network")
	if err != nil {
		return false, err
	}
	for _, container := range x.Containers {
		if container.Name == "amazeeio-haproxy" {
			return true, nil
		}
	}
	return false, nil
}
