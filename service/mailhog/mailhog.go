package mailhog

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	model "github.com/fubarhouse/pygmy/v1/service/interface"
)

func New() model.Service {
	return model.Service{
		ContainerName: "mailhog.docker.amazee.io",
		Config:        container.Config{
			User:       "0",
			ExposedPorts: nat.PortSet{
				"80/tcp": struct{}{},
				"1025/tcp": struct{}{},
				"8025/tcp": struct{}{},
			},
			Env: []string{
				"MH_UI_BIND_ADDR=0.0.0.0:80",
				"MH_API_BIND_ADDR=0.0.0.0:80",
				"AMAZEEIO=AMAZEEIO",
			},
			Image: "mailhog/mailhog",
			Labels:		map[string]string{
				"pygmy": "pygmy",
			},
		},
		HostConfig:    container.HostConfig{
			AutoRemove: false,
			RestartPolicy: struct {
				Name              string
				MaximumRetryCount int
			}{Name: "always", MaximumRetryCount: 0},
			PortBindings: nat.PortMap{
				"1025/tcp": []nat.PortBinding{
					{
						HostIP: "",
						HostPort: "1025",
					},
				},
			},
		},
		NetworkConfig: network.NetworkingConfig{},
	}

}
