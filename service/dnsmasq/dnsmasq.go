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
			"6053:53/tcp",
			"-p",
			"6053:53/udp",
			"--name=amazeeio-dnsmasq",
			"--cap-add=NET_ADMIN",
			"andyshinn/dnsmasq:2.78",
			"-A",
			"/docker.amazee.io/127.0.0.1",
		},
	}
}
