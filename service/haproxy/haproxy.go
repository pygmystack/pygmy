package haproxy

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy/v1/service/interface"
)

func New() model.Service {
	return model.Service{
		ContainerName: "amazeeio-haproxy",
		Config:        container.Config{
			Image:    "amazeeio/haproxy",
			Labels:		map[string]string{
				"pygmy": "pygmy",
			},

		},
		HostConfig:    container.HostConfig{
			Binds: []string{"/var/run/docker.sock:/tmp/docker.sock"},
			AutoRemove: false,
			RestartPolicy: struct {
				Name              string
				MaximumRetryCount int
			}{Name: "always", MaximumRetryCount: 0},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}