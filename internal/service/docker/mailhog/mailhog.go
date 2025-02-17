package mailhog

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	"github.com/pygmystack/pygmy/internal/runtime/docker"
)

// New will provide the standard object for the mailhog container.
func New(c *docker.Params) docker.Service {
	return docker.Service{
		Config: container.Config{
			User: "0",
			ExposedPorts: nat.PortSet{
				"80/tcp":   struct{}{},
				"1025/tcp": struct{}{},
				"8025/tcp": struct{}{},
			},
			Env: []string{
				"MH_UI_BIND_ADDR=0.0.0.0:80",
				"MH_API_BIND_ADDR=0.0.0.0:80",
				"AMAZEEIO=AMAZEEIO",
				fmt.Sprintf("AMAZEEIO_URL=mailhog.%s", c.Domain),
			},
			Image: "pygmystack/mailhog",
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     "amazeeio-mailhog",
				"pygmy.network":  "amazeeio-network",
				"pygmy.url":      fmt.Sprintf("http://mailhog.%s", c.Domain),
				"pygmy.weight":   "15",
			},
		},
		HostConfig: container.HostConfig{
			AutoRemove: false,
			RestartPolicy: container.RestartPolicy{
				Name:              "unless-stopped",
				MaximumRetryCount: 0,
			},
		},
		NetworkConfig: network.NetworkingConfig{},
	}

}

// NewDefaultPorts will provide the standard ports used for merging into the
// mailhog config.
func NewDefaultPorts() docker.Service {
	return docker.Service{
		HostConfig: container.HostConfig{
			PortBindings: nat.PortMap{
				"1025/tcp": []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: "1025",
					},
				},
			},
		},
	}
}
