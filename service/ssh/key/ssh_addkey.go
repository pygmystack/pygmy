// +build darwin linux

package key

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy/service/interface"
)

func NewAdder(key string) model.Service {
	return model.Service{
		Name:     "amazeeio-ssh-agent-add-key",
		Weight:   31,
		Discrete: true,
		Output:   true,
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Cmd: []string{
				"ssh-add",
				key,
			},
			Labels: map[string]string{
				"pygmy": "pygmy",
			},
		},
		HostConfig: container.HostConfig{
			IpcMode:     "private",
			AutoRemove:  true,
			Binds:       []string{fmt.Sprintf("%v:%v", key, key)},
			VolumesFrom: []string{"amazeeio-ssh-agent"},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}

func NewShower() model.Service {
	return model.Service{
		Name:     "amazeeio-ssh-agent-show-keys",
		Weight:   32,
		Discrete: true,
		Output:   true,
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Cmd: []string{
				"ssh-add",
				"-l",
			},
			Labels: map[string]string{
				"pygmy": "pygmy",
			},
		},
		HostConfig: container.HostConfig{
			AutoRemove:  true,
			VolumesFrom: []string{"amazeeio-ssh-agent"},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
