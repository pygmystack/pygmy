package ssh_agent

import (
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy/service/interface"
	"github.com/fubarhouse/pygmy/service/ssh/ssh_addkey"
)

func New() model.Service {
	return model.Service{
		ContainerName: "amazeeio-ssh-agent",
		Config:        container.Config{
			Image:    "amazeeio/ssh-agent",
			Labels:		map[string]string{
				"pygmy": "pygmy",
			},
		},
		HostConfig:    container.HostConfig{
			AutoRemove:  false,
			RestartPolicy: struct {
				Name              string
				MaximumRetryCount int
			}{Name: "always", MaximumRetryCount: 0},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}

func List() ([]byte, error) {
	i := ssh_addkey.NewShower()
	return i.Start()
}

func Search(key string) bool {
	items, _ := List()
	for _, item := range strings.Split(string(items), "\n") {
		if strings.Contains(item, "The agent has no identities") {
			return false
		}
		if strings.Contains(item, key) {
			return true
		}
	}
	return false
}