package agent

import (
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy/v1/service/interface"
	"github.com/fubarhouse/pygmy/v1/service/ssh/key"
)

func New() model.Service {
	return model.Service{
		Name: "amazeeio-ssh-agent",
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

func List() []byte {
	i := key.NewShower()
	r, _ := i.Start()
	return r
}

func Search(key string) bool {
	items := List()
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