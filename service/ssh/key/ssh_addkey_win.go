// +build windows

package key

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy/v1/service/interface"
)

func NewAdder(key string) model.Service {
	return model.Service{
		Name:     "amazeeio-ssh-agent-add-key",
		Discrete: true,
		Output:   true,
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Cmd: []string{
				"windows-key-add",
				"/key",
			},
			Labels: map[string]string{
				"pygmy": "pygmy",
			},
		},
		HostConfig: container.HostConfig{
			IpcMode:     "private",
			AutoRemove:  true,
			Binds:       []string{fmt.Sprintf("%v:/%v", key, key)},
			VolumesFrom: []string{"amazeeio-ssh-agent"},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}

func NewShower() model.Service {
	return model.Service{
		Name: "amazeeio-ssh-agent-show-keys",
		Discrete: true,
		Output: true,
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Cmd: []string{
				"ssh-add",
				"-l",
			},
			Labels:		map[string]string{
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
