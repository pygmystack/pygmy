// +build windows

package ssh_addkey

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy/service/interface"
	"github.com/mitchellh/go-homedir"
	"os"
)

func NewAdder(key string) model.Service {
	if key == "" {
		homedir, _ := homedir.Dir()
		key = fmt.Sprintf("%v%v.ssh%vid_rsa", homedir, string(os.PathSeparator), string(os.PathSeparator))
	}

	return model.Service{
		ContainerName: "amazeeio-ssh-agent-add-key",
		Config:        container.Config{
			Image:    "amazeeio/ssh-agent",
			Cmd: []string{
				"windows-key-add",
				"/key",
			},
		},
		HostConfig:    container.HostConfig{
			AutoRemove:  true,
			Mounts: []mount.Mount{
				{
					Type: "bind",
					Source: key,
					Target: key,
					ReadOnly: true,
				},
			},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}

func NewShower() model.Service {
	return model.Service{
		ContainerName: "amazeeio-ssh-agent-add-key",
		Config:        container.Config{
			Image:    "amazeeio/ssh-agent",
			Cmd: []string{
				"ssh-add",
				"-l",
			},
		},
		HostConfig:    container.HostConfig{
			AutoRemove:  true,
			VolumesFrom: []string{"amazeeio-ssh-agent"},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
