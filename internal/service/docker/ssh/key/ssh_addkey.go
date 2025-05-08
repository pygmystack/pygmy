//go:build !windows
// +build !windows

package key

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"

	"github.com/pygmystack/pygmy/internal/runtime/docker"
)

// NewAdder will provide the standard object for the SSH key adder container.
func NewAdder() docker.Service {
	return docker.Service{
		Config: container.Config{
			Image: "pygmystack/ssh-agent",
			Labels: map[string]string{
				"pygmy.defaults":    "true",
				"pygmy.enable":      "true",
				"pygmy.name":        "amazeeio-ssh-agent-add-key",
				"pygmy.network":     "amazeeio-network",
				"pygmy.discrete":    "true",
				"pygmy.interactive": "true",
				"pygmy.output":      "false",
				"pygmy.purpose":     "addkeys",
				"pygmy.weight":      "31",
			},
			Tty:       true,
			OpenStdin: true,
		},
		HostConfig: container.HostConfig{
			AutoRemove:  true,
			IpcMode:     "private",
			VolumesFrom: []string{"amazeeio-ssh-agent"},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}
