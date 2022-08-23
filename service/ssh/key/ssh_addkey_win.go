//go:build windows
// +build windows

package key

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/pygmystack/pygmy/service/interface"
)

// NewAdder will provide the standard object for the SSH key adder container.
func NewAdder(c *model.Params) model.Service {
	return model.Service{
		Config: container.Config{
			Image: "pygmystack/ssh-agent",
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     fmt.Sprintf("%s-ssh-agent-add-key", c.Prefix),
				"pygmy.network":  fmt.Sprintf("%s-network", c.Prefix),
				"pygmy.discrete": "true",
				"pygmy.output":   "false",
				"pygmy.purpose":  "addkeys",
				"pygmy.weight":   "31",
			},
		},
		HostConfig: container.HostConfig{
			IpcMode:     "private",
			AutoRemove:  false,
			VolumesFrom: []string{fmt.Sprintf("%s-ssh-agent", c.Prefix)},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
