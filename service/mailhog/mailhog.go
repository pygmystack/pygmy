// Mailhog provides default values for the Mailhog docker container.
package mailhog

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	model "github.com/fubarhouse/pygmy-go/service/model"
)

func New() model.Service {
	return model.Service{
		Name: "mailhog.docker.amazee.io",
		URL: "http://mailhog.docker.amazee.io",
		Weight: 15,
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
		},
		NetworkConfig: network.NetworkingConfig{},
	}

}

// NewDefaultPorts returns a bare struct containing the PortBindings for
// the container. This will be merged into the service if PortBindings
// are missing from the struct which Pygmy resolves.
func NewDefaultPorts() model.Service {
	return model.Service{
		HostConfig: container.HostConfig{
			PortBindings: nat.PortMap{
				"1025/tcp": []nat.PortBinding{
					{
						HostIP: "",
						HostPort: "1025",
					},
				},
			},
		},
	}
}