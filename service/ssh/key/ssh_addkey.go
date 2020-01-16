// +build darwin linux

package key

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy-go/service/interface"
)

func NewAdder() model.Service {
	return model.Service{
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Labels: map[string]string{
				"pygmy":          "pygmy",
				"pygmy.name":     "amazeeio-ssh-agent-add-key",
				"pygmy.discrete": "true",
				"pygmy.output":   "true",
				"pygmy.purpose":  "addkeys",
				"pygmy.weight":   "31",
			},
		},
		HostConfig: container.HostConfig{
			IpcMode:     "private",
			AutoRemove:  true,
			VolumesFrom: []string{"amazeeio-ssh-agent"},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}

func NewShower() model.Service {
	return model.Service{
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Cmd: []string{
				"ssh-add",
				"-l",
			},
			Labels: map[string]string{
				"pygmy":          "pygmy",
				"pygmy.name":     "amazeeio-ssh-agent-show-keys",
				"pygmy.discrete": "true",
				"pygmy.output":   "false",
				"pygmy.purpose":  "showkeys",
				"pygmy.weight":   "32",
			},
		},
		HostConfig: container.HostConfig{
			AutoRemove:  true,
			VolumesFrom: []string{"amazeeio-ssh-agent"},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
