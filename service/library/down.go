package library

import (
	"fmt"
	"github.com/fubarhouse/pygmy-go/service/interface/docker"

	. "github.com/logrusorgru/aurora"
)

// Down will bring pygmy down safely
func Down(c Config) {

	Setup(&c)
	for _, Service := range c.Services {
		enabled, _ := Service.GetFieldBool("enable")
		name, _ := Service.GetFieldString("name")
		if enabled {
			e := Service.Stop()
			if e != nil {
				fmt.Print(Red(fmt.Sprintf("Successfully removed %s\n", name)))
			}
		}
	}

	for _, network := range c.Networks {
		NetworksToClean = append(NetworksToClean, network.Name)
	}

	for _, network := range unique(NetworksToClean) {
		e := docker.DockerNetworkRemove(network)
		if e != nil {
			fmt.Println(e)
		}
		if s, _ := docker.DockerNetworkStatus(network); !s {
			color.Print(Green(fmt.Sprintf("Successfully removed network %s\n", network)))
		} else {
			color.Print(Red(fmt.Sprintf("Network %s was not removed\n", network)))
		}
	}

	for _, resolver := range c.Resolvers {
		resolver.Clean()
	}
}
