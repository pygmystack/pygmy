package dnsmasq

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	model "github.com/fubarhouse/pygmy-go/service/interface"
)

func New() model.Service {
	return model.Service{
		Name:   "amazeeio-dnsmasq",
		Weight: 13,
		Config: container.Config{
			Image: "andyshinn/dnsmasq:2.78",
			Cmd: []string{
				"-A",
				"/docker.amazee.io/127.0.0.1",
			},
			Labels: map[string]string{
				"pygmy":        "pygmy",
				"pygmy.name":   "amazeeio-dnsmasq",
				"pygmy.weight": "13",
			},
		},
		HostConfig: container.HostConfig{
			AutoRemove: false,
			CapAdd:     []string{"NET_ADMIN"},
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
			}{Name: "no", MaximumRetryCount: 0},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
