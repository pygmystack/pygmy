// HAProxy provides default values for the HAProxy docker container.
package haproxy

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	model "github.com/fubarhouse/pygmy-go/service/model"
)

// New will return a data structure containing all the default values
// Pygmy needs for the HAProxy docker container. It does not contain
// the port configuration - if the port configuration is missing, it
// will merge in the response struct from NewDefaultPorts().
func New() model.Service {
	return model.Service{
		Name:   "amazeeio-haproxy",
		URL:    "http://docker.amazee.io/stats",
		Weight: 14,
		Config: container.Config{
			Image: "amazeeio/haproxy",
			Labels: map[string]string{
				"pygmy": "pygmy",
			},
		},
		HostConfig: container.HostConfig{
			Binds:        []string{"/var/run/docker.sock:/tmp/docker.sock"},
			AutoRemove:   false,
			PortBindings: nil,
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
