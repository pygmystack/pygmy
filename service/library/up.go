package library

import (
	"fmt"
	"github.com/fubarhouse/pygmy/v1/service/haproxy_connector"
	model "github.com/fubarhouse/pygmy/v1/service/interface"
	"github.com/fubarhouse/pygmy/v1/service/network"
	"github.com/fubarhouse/pygmy/v1/service/resolv"
)

func Up(c Config) {

	Setup(&c)

	for _, volume := range c.Volumes {
		if s, _ := model.DockerVolumeExists(volume); !s {
			_, err := model.DockerVolumeCreate(volume)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("Already created volume %v\n", volume)
		}
	}

	// Maps are... bad for predictable sequencing.
	// Look over the sorted slice and start them in
	// alphabetical order - so that one can configure
	// an ssh-agent like amazeeio-ssh-agent.
	for _, service := range c.SortedServices {
		s := c.Services[service]
		if !s.Disabled {
			s.Start()
		}
	}

	for Network, Containers := range c.Networks {
		netStat, _ := network.Status(Network)
		if !netStat {
			network.Create(Network)
		}
		for _, Container := range Containers {
			if s, _ := haproxy_connector.Connected(Container, Network); !s {
				haproxy_connector.Connect(Container, Network)
				if s, _ := haproxy_connector.Connected(Container, Network); s {
					fmt.Printf("Successfully connected %v to %v\n", Container, Network)
				} else {
					fmt.Printf("Could not connect %v to %v\n", Container, Network)
				}
			} else {
				fmt.Printf("Already connected %v to %v\n", Container, Network)
			}
		}
	}

	if !c.SkipResolver {
		for _, resolver := range c.Resolvers {
			resolv.New(resolv.Resolv{Name: resolver.Name, Data: resolver.Data, Folder: resolver.Folder, File: resolver.File}).Configure()
		}
	}

	if !c.SkipKey {

		SshKeyAdd(c, c.Key)

	}
}
