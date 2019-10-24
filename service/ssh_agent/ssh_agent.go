package ssh_agent

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy/service/interface"
)

func New() model.Service {
	return model.Service{
		ContainerName: "amazeeio-ssh-agent",
		Config:        container.Config{
			Image:    "amazeeio/ssh-agent",
		},
		HostConfig:    container.HostConfig{
			AutoRemove:  false,
			RestartPolicy: struct {
				Name              string
				MaximumRetryCount int
			}{Name: "always", MaximumRetryCount: 0},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}

