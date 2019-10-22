package haproxy

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	model "github.com/fubarhouse/pygmy/service/interface"
)

func New() model.Service {
	return model.Service{
		ContainerName: "amazeeio-haproxy",
		Config:        container.Config{
			Image:    "amazeeio/haproxy",

		},
		HostConfig:    container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type: "bind",
					Source: "/var/run/docker.sock",
					Target: "/tmp/docker.sock",
					ReadOnly: false,
				},
			},
			AutoRemove: false,
			RestartPolicy: struct {
				Name              string
				MaximumRetryCount int
			}{Name: "always", MaximumRetryCount: 0},
			PortBindings: nat.PortMap{
				"80/tcp": []nat.PortBinding{
					{
						HostIP: "0.0.0.0",
						HostPort: "80",
					},
				},
				"443/tcp": []nat.PortBinding{
					{
						HostIP: "0.0.0.0",
						HostPort: "443",
					},
				},
			},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}