package dnsmasq

import (
	model "github.com/fubarhouse/pygmy/service/interface"
)

func New() model.Service {
	return model.Service{
		Name:          "Dnsmasq",
		Address:       "127.0.0.1",
		ContainerName: "amazeeio-dnsmasq",
		Domain:        "docker.amazee.io",
		ImageName:     "andyshinn/dnsmasq:2.75",
		RunCmd: []string{
			"run",
			"-d",
			"-p",
			"53:53/tcp",
			"-p",
			"53:53/udp",
			"--name=amazeeio-dnsmasq",
			"--cap-add=NET_ADMIN",
			"andyshinn/dnsmasq:2.75",
			"-A",
			"/docker.amazee.io/127.0.0.1",
		},
	}
}
