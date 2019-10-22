package network

import (
	model "github.com/fubarhouse/pygmy/service/interface"
	"strings"
)

func Create() error {
	return model.DockerNetworkCreate("amazeeio-network")
}

func Status() (bool, error) {
	output, error := model.DockerRun([]string{"network", "ls", "--format", "'{{.Name}}'"})
	if strings.Contains(string(output), "amazeeio-network") {
		return true, nil
	}
	return false, error
}
