package library

import (
	"fmt"

	model "github.com/fubarhouse/pygmy-go/service/interface"
)

// Down will bring pygmy down safely
func Down(c Config) {

	Setup(&c)

	for _, Service := range c.Services {
		enabled, _ := Service.GetFieldBool("enable")
		if enabled {
			Service.Stop()
		}
	}

	for _, network := range c.Networks {
		model.DockerNetworkRemove(&network)
		if s, _ := model.DockerNetworkStatus(&network); s {
			fmt.Printf("Successfully removed network %v\n", network.Name)
		} else {
			fmt.Printf("Network %v was not removed", network.Name)
		}
	}

	for _, resolver := range c.Resolvers {
		resolver.Clean()
	}
}
