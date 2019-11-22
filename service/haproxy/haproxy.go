package haproxy

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy/service/interface"
)

func New() model.Service {
	return model.Service{
		Name: "amazeeio-haproxy",
		Weight: 14,
		Config: container.Config{
			Image: "amazeeio/haproxy",
			Labels: map[string]string{
				"pygmy": "pygmy",
			},
		},
		HostConfig: container.HostConfig{
			Binds:      []string{"/var/run/docker.sock:/tmp/docker.sock"},
			AutoRemove: false,
			RestartPolicy: struct {
				Name              string
				MaximumRetryCount int
			}{Name: "always", MaximumRetryCount: 0},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
