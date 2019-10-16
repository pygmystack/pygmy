package haproxy_connector

import (
	model "github.com/fubarhouse/pygmy/service/interface"
	"strings"
)

func Connect() error {
	if s, _ := Connected(); !s {
		_, error := model.DockerRun([]string{"network", "connect", "amazeeio-network", "amazeeio-haproxy"})
		return error
	}
	return nil
}

func Connected() (bool, error) {
	output, error := model.DockerRun([]string{"network", "inspect", "amazeeio-network", "-f", "'{{.Containers}}'"})
	if strings.Contains(string(output), "amazeeio-haproxy") {
		return true, nil
	}
	return false, error
}
