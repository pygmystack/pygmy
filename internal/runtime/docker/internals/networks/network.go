package networks

import (
	"context"
	"fmt"

	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/containers"
)

// Create is an abstraction layer on top of the Docker API call
// which will create a Docker network using a specified configuration.
func Create(ctx context.Context, cli *client.Client, network *networktypes.Inspect) error {
	netVal, _ := Status(ctx, cli, network.Name)
	if netVal {
		return fmt.Errorf("docker network %v already exists", network.Name)
	}

	config := networktypes.CreateOptions{
		Driver:     network.Driver,
		EnableIPv6: &network.EnableIPv6,
		IPAM:       &network.IPAM,
		Internal:   network.Internal,
		Attachable: network.Attachable,
		Options:    network.Options,
		Labels:     network.Labels,
	}
	_, err := cli.NetworkCreate(ctx, network.Name, config)
	if err != nil {
		return err
	}

	return nil
}

// Remove will attempt to remove a Docker network
// and will not apply force to removal.
func Remove(ctx context.Context, cli *client.Client, network string) error {
	err := cli.NetworkRemove(ctx, network)
	if err != nil {
		return err
	}
	return nil
}

// Status will identify if a network with a
// specified name is present been created and return a boolean.
func Status(ctx context.Context, cli *client.Client, network string) (bool, error) {
	networks, err := cli.NetworkList(ctx, networktypes.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, n := range networks {
		if n.Name == network {
			return true, nil
		}
	}

	return false, nil
}

// Get will use the Docker API to retrieve a Docker network
// which has a given name.
func Get(ctx context.Context, cli *client.Client, name string) (networktypes.Inspect, error) {
	networks, err := cli.NetworkList(ctx, networktypes.ListOptions{})
	if err != nil {
		return networktypes.Inspect{}, err
	}
	for _, network := range networks {
		if val, ok := network.Labels["pygmy.name"]; ok {
			if val == name {
				return network, nil
			}
		}
	}
	return networktypes.Inspect{}, nil
}

// Connect will connect a container to a network.
func Connect(ctx context.Context, cli *client.Client, network string, containerName string) error {
	e := cli.NetworkConnect(ctx, network, containerName, nil)
	if e != nil {
		return e
	}
	return nil
}

// Connected will check if a container is connected to a network.
func Connected(ctx context.Context, cli *client.Client, network string, containerName string) (bool, error) {
	// Reset network state:
	c, _ := containers.List(ctx, cli)
	for d := range c {
		if c[d].Labels["pygmy.name"] == containerName {
			for net := range c[d].NetworkSettings.Networks {
				if net == network {
					return true, nil
				}
			}
		}
	}
	return false, fmt.Errorf("network was found without the container connected")
}
