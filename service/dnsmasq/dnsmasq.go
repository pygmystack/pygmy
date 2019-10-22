package dnsmasq

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	model "github.com/fubarhouse/pygmy/service/interface"
)

func New() model.Service {
	return model.Service{
		ContainerName: "amazeeio-dnsmasq",
		Config:        container.Config{
			Image: "andyshinn/dnsmasq:2.78",
			Cmd: []string{
				"-A",
				"/docker.amazee.io/127.0.0.1",
			},},
		HostConfig:    container.HostConfig{
			AutoRemove: false,
			CapAdd:     []string{"NET_ADMIN"},
			PortBindings: nat.PortMap{
				"6053/tcp": []nat.PortBinding{
					{
						HostIP: "0.0.0.0",
						HostPort: "6053",
					},
				},
				"6053/udp": []nat.PortBinding{
					{
						HostIP: "0.0.0.0",
						HostPort: "6053",
					},
				},
			},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
