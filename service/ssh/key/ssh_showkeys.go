// +build darwin linux

package key

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy-go/service/interface"
)

// NewShower will provide the standard object for the SSH key shower container.
func NewShower() model.Service {
	return model.Service{
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Cmd: []string{
				"ssh-add",
				"-L",
			},
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     "amazeeio-ssh-agent-show-keys",
				"pygmy.network":  "amazeeio-network",
				"pygmy.discrete": "true",
				"pygmy.output":   "false",
				"pygmy.purpose":  "showkeys",
				"pygmy.weight":   "32",
			},
		},
		HostConfig: container.HostConfig{
			AutoRemove:  true,
			IpcMode:     "private",
			VolumesFrom: []string{"amazeeio-ssh-agent"},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
