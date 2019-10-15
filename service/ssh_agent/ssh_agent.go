package ssh_agent

import (
	model "github.com/fubarhouse/pygmy/service/interface"
)

func New() model.Service {
	return model.Service{
		Name: "amazeeio-ssh-agent",
		Address: "",
		ContainerName: "amazeeio-ssh-agent",
		Domain: "",
		ImageName: "amazeeio/ssh-agent",
		RunCmd: []string{
			"run",
			"-d",
			"--restart=always",
			"--name",
			"amazeeio-ssh-agent",
			"amazeeio/ssh-agent",
		},

	}
}

