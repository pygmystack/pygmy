package haproxy

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	model "github.com/pygmystack/pygmy/service/interface"
	"github.com/pygmystack/pygmy/service/interface/docker"
	"strings"
)

// New will provide the standard object for the haproxy container.
func New(c *model.Params) model.Service {
	socket := strings.TrimLeft(docker.GetSocket(), "unix:///")
	return model.Service{
		Config: container.Config{
			Image: "pygmystack/haproxy",
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     "amazeeio-haproxy",
				"pygmy.network":  "amazeeio-network",
				"pygmy.url":      fmt.Sprintf("http://%s/stats", c.Domain),
				"pygmy.weight":   "14",
			},
			Env: []string{
				fmt.Sprintf("AMAZEEIO_URL=%s", c.Domain),
			},
		},
		HostConfig: container.HostConfig{
			Binds:        []string{fmt.Sprintf("/%s:/tmp/docker.sock", socket)},
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
