package haproxy

import (
	model "github.com/fubarhouse/pygmy/service/interface"
)

func New() model.Service {
	return model.Service{
		Name:          "Haproxy",
		Address: "",
		ContainerName: "amazeeio-haproxy",
		Domain: "",
		ImageName: "amazeeio/haproxy",
		RunCmd: []string{
			"run",
			"-d",
			"-p",
			"80:80",
			"--volume=/var/run/docker.sock:/tmp/docker.sock",
			"--restart=always",
			"--name=amazeeio-haproxy",
			"amazeeio/haproxy",
		},
	}
}

