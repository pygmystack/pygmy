package haproxy

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	"github.com/pygmystack/pygmy/internal/runtime/docker"
)

// New will provide the standard object for the haproxy container.
func New(c *docker.Params, tlsCertPath string) docker.Service {
	binds := []string{"/var/run/docker.sock:/tmp/docker.sock"}
	if tlsCertPath != "" {
		binds = append(binds, fmt.Sprintf("%s:/app/server.pem:ro", tlsCertPath))
	}
	return docker.Service{
		Config: container.Config{
			Image: "pygmystack/haproxy",
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     "amazeeio-haproxy",
				"pygmy.network":  "amazeeio-network",
				"pygmy.url":      fmt.Sprintf("%s/stats", c.Domain),
				"pygmy.weight":   "14",
			},
			Env: []string{
				fmt.Sprintf("AMAZEEIO_URL=%s", c.Domain),
			},
		},
		HostConfig: container.HostConfig{
			Binds:        binds,
			AutoRemove:   false,
			PortBindings: nil,
			RestartPolicy: container.RestartPolicy{
				Name:              "unless-stopped",
				MaximumRetryCount: 0,
			},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}

// NewDefaultPorts will provide the standard ports used for merging into the
// haproxy config.
func NewDefaultPorts() docker.Service {
	return docker.Service{
		HostConfig: container.HostConfig{
			PortBindings: nat.PortMap{
				"80/tcp": []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: "80",
					},
				},
				"443/tcp": []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: "443",
					},
				},
			},
		},
	}
}
