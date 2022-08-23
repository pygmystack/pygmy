package dnsmasq

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	model "github.com/pygmystack/pygmy/service/interface"
)

// New will provide the standard object for the dnsmasq container.
func New(c *model.Params) model.Service {
	return model.Service{
		Config: container.Config{
			Image: "pygmystack/dnsmasq",
			Cmd: []string{
				"--log-facility=-",
				"-A",
				fmt.Sprintf("/%s/127.0.0.1", c.Domain),
			},
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     fmt.Sprintf("%s-dnsmasq", c.Prefix),
				"pygmy.weight":   "13",
			},
		},
		HostConfig: container.HostConfig{
			AutoRemove: false,
			CapAdd:     []string{"NET_ADMIN"},
			IpcMode:    "private",
			PortBindings: nat.PortMap{
				"53/tcp": []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: "6053",
					},
				},
				"53/udp": []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: "6053",
					},
				},
			},
			RestartPolicy: struct {
				Name              string
				MaximumRetryCount int
			}{Name: "unless-stopped", MaximumRetryCount: 0},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
