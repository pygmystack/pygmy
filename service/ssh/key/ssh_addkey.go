// +build darwin linux

// Key provides default values for the SSH Key adder and shower docker container(s).
package key

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy-go/service/model"
)

func NewAdder() model.Service {
	return model.Service{
		Name:     "amazeeio-ssh-agent-add-key",
		Group:    "addkeys",
		Weight:   31,
		Discrete: true,
		Output:   true,
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Labels: map[string]string{
				"pygmy": "pygmy",
				"pygmy.addkeys": "pygmy.addkeys",
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
		Name:     "amazeeio-ssh-agent-show-keys",
		Group:    "showkeys",
		Weight:   32,
		Discrete: true,
		Output:   false,
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Cmd: []string{
				"ssh-add",
				"-l",
			},
			Labels: map[string]string{
				"pygmy": "pygmy",
				"pygmy.showkeys": "pygmy.showkeys",
			},
		},
		HostConfig: container.HostConfig{
			AutoRemove:  true,
			VolumesFrom: []string{"amazeeio-ssh-agent"},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
