package agent

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	model "github.com/pygmystack/pygmy/service/interface"
)

// New will provide the standard object for the SSH agent container.
func New(c *model.Params) model.Service {
	return model.Service{
		Config: container.Config{
			Image: "pygmystack/ssh-agent",
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     fmt.Sprintf("%s-ssh-agent", c.Prefix),
				"pygmy.network":  fmt.Sprintf("%s-network", c.Prefix),
				"pygmy.output":   "false",
				"pygmy.purpose":  "sshagent",
				"pygmy.weight":   "10",
			},
		},
		HostConfig: container.HostConfig{
			AutoRemove: false,
			IpcMode:    "private",
			RestartPolicy: struct {
				Name              string
				MaximumRetryCount int
			}{Name: "unless-stopped", MaximumRetryCount: 0},
		},
		NetworkConfig: network.NetworkingConfig{},
	}
}

// List will grab the output of all running containers with the proper
// config after starting them, and return it.
// which is indicated by the purpose tag.
func List(service model.Service) ([]byte, error) {
	purpose, _ := service.GetFieldString("purpose")
	if purpose == "showkeys" {
		e := service.Start()
		if e != nil {
			return []byte{}, e
		}
	}
	return service.DockerLogs()
}

// Validate will validate if an SSH key is valid.
func Validate(filePath, passcode string) (bool, error) {

	filePath = strings.TrimRight(filePath, ".pub")
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Err")
	}

	if passcode == "" {
		_, err = ssh.ParsePrivateKey(content)
		if err != nil {
			return false, err
		}
	} else {
		_, err = ssh.ParsePrivateKeyWithPassphrase(content, []byte(passcode))
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// Search will determine if an SSH key has been added to the agent.
func Search(service model.Service, key string) (bool, error) {
	result := false
	if _, err := os.Stat(key); !os.IsNotExist(err) {
		stripped := strings.Trim(key, ".pub")
		data, err := ioutil.ReadFile(stripped + ".pub")
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
