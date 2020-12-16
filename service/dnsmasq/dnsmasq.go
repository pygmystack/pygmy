package dnsmasq

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	model "github.com/fubarhouse/pygmy-go/service/interface"
)

// New will provide the standard object for the dnsmasq container.
func New(domain string) model.Service {
	return model.Service{
		Config: container.Config{
			Image: "andyshinn/dnsmasq:2.78",
			Cmd: []string{
				"-A",
				fmt.Sprintf("/%v/127.0.0.1", domain),
			},
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     "amazeeio-dnsmasq",
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
