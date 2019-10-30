// +build darwin linux

package ssh_addkey

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy/service/interface"
)

func NewAdder(key string) model.Service {
	return model.Service{
		ContainerName: "amazeeio-ssh-agent-add-key",
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Cmd: []string{
				"ssh-add",
				key,
			},
			Labels:		map[string]string{
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
		ContainerName: "amazeeio-ssh-agent-show-keys",
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
