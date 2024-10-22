package agent

import (
	"errors"
	"fmt"
	"github.com/pygmystack/pygmy/internal/runtime"
	"github.com/pygmystack/pygmy/internal/runtime/docker"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/containers"
)

// New will provide the standard object for the SSH agent container.
func New() docker.Service {
	return docker.Service{
		Config: container.Config{
			Image: "pygmystack/ssh-agent",
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     "amazeeio-ssh-agent",
				"pygmy.network":  "amazeeio-network",
				"pygmy.output":   "false",
				"pygmy.purpose":  "sshagent",
				"pygmy.weight":   "10",
			},
		},
		HostConfig: container.HostConfig{
			AutoRemove: false,
			IpcMode:    "private",
			RestartPolicy: container.RestartPolicy{
				Name:              "unless-stopped",
				MaximumRetryCount: 0,
			},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}

// List will grab the output of all running containers with the proper
// config after starting them, and return it.
// which is indicated by the purpose tag.
func List(service *runtime.Service) ([]byte, error) {
	name, _ := service.GetFieldString("name")
	purpose, _ := service.GetFieldString("purpose")
	if purpose == "showkeys" {
		e := service.Start()
		if e != nil {
			return []byte{}, e
		}
	}
	return containers.Logs(name)
}

// Validate will validate if an SSH key is valid.
func Validate(filePath string) (bool, error) {

	filePath = strings.TrimRight(filePath, ".pub")
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Err")
	}

	_, err = ssh.ParsePrivateKey(content)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Search will determine if an SSH key has been added to the agent.
func Search(service *runtime.Service, key string) (bool, error) {
	result := false
	if _, err := os.Stat(key); !os.IsNotExist(err) {
		stripped := strings.Trim(key, ".pub")
		data, err := os.ReadFile(stripped + ".pub")
		if err != nil {
			return false, err
		}

		items, _ := List(service)

		if len(items) == 0 {
			return false, nil
		}

		for _, item := range strings.Split(string(items), "\n") {
			if strings.Contains(item, "The agent has no identities") {
				return false, errors.New(item)
			}
			if strings.Contains(item, "Error loading key") {
				return false, errors.New(item)
			}
			if strings.Contains(item, string(data)) {
				result = true
			}
		}
	}
	return result, nil
}