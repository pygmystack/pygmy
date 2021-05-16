package haproxy

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	model "github.com/fubarhouse/pygmy-go/service/interface"
)

// New will provide the standard object for the haproxy container.
func New() model.Service {
	return model.Service{
		Config: container.Config{
			Image: "amazeeio/haproxy",
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     "amazeeio-haproxy",
				"pygmy.network":  "amazeeio-network",
				"pygmy.url":      "http://docker.amazee.io/stats",
				"pygmy.weight":   "14",
			},
			Healthcheck: &container.HealthConfig{
				Test:        []string{"CMD-SHELL", "wget http://docker.amazee.io/stats -q -S -O - 2>&1 | grep docker.amazee.io"},
				Interval:    30000000000,
				Timeout:     5000000000,
				StartPeriod: 5000000000,
				Retries:     5,
			},
		},
		HostConfig: container.HostConfig{
			Binds:        []string{"/var/run/docker.sock:/tmp/docker.sock"},
			AutoRemove:   false,
			PortBindings: nil,
			RestartPolicy: struct {
				Name              string
				MaximumRetryCount int
			}{Name: "unless-stopped", MaximumRetryCount: 0},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}

// NewDefaultPorts will provide the standard ports used for merging into the
// haproxy config.
func NewDefaultPorts() model.Service {
	return model.Service{
		HostConfig: container.HostConfig{
			PortBindings: nat.PortMap{
				"80/tcp": []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: "80",
					},
				},
			},
		},
	}
}
