package agent

import (
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy-go/service/interface"
)

// New will provide the standard object for the SSH agent container.
func New() model.Service {
	return model.Service{
		Config: container.Config{
			Image: "amazeeio/ssh-agent",
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     "amazeeio-ssh-agent",
				"pygmy.output":   "false",
				"pygmy.purpose":  "sshagent",
				"pygmy.weight":   "30",
			},
		},
		HostConfig: container.HostConfig{
			AutoRemove: false,
			IpcMode:    "private",
			RestartPolicy: struct {
				Name              string
				MaximumRetryCount int
			}{Name: "on-failure", MaximumRetryCount: 0},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}

// SshKeyLister will grab the output of all running containers with the proper
// config after starting them, and return it.
// which is indicated by the purpose tag.
func List(service model.Service) []string {
	var r []byte
	purpose, _ := service.GetFieldString("purpose")
	if purpose == "showkeys" {
		r, _ = service.Start()
	}
	return strings.Split(string(r), "\n")
}

// Search will determine if an SSH key has been added to the agent.
func Search(service model.Service, key string) bool {
	if _, err := os.Stat(key); !os.IsNotExist(err) {

		items := List(service)

		if len(items) == 0 {
			return false
		}

		for _, item := range items {
			if strings.Contains(item, "The agent has no identities") {
				return false
			}
			if strings.Contains(item, key) {
				return true
			}
		}
	}
	return false
}
