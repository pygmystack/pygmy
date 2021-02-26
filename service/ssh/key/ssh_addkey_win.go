// +build windows

package key

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy-go/service/interface"
)

// NewAdder will provide the standard object for the SSH key adder container.
func NewAdder() model.Service {
	return model.Service{
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     "amazeeio-ssh-agent-add-key",
				"pygmy.network":  "amazeeio-network",
				"pygmy.discrete": "true",
				"pygmy.output":   "false",
				"pygmy.purpose":  "addkeys",
				"pygmy.weight":   "31",
			},
		},
		HostConfig: container.HostConfig{
			IpcMode:     "private",
			AutoRemove:  false,
			VolumesFrom: []string{"amazeeio-ssh-agent"},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
