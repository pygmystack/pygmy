package library

import (
	"fmt"

	model "github.com/fubarhouse/pygmy-go/service/interface"
)

// Clean will forcibly kill and remove all of pygmy's containers in the daemon
func Clean(c Config) {

	Setup(&c)
	Containers, _ := model.DockerContainerList()

	for _, Container := range Containers {
		target := false
		if l := Container.Labels["pygmy.enable"]; l == "true" || l == "1" {
			target = true
		}
		if l := Container.Labels["pygmy"]; l == "pygmy" {
			target = true
		}

		if target {
			err := model.DockerKill(Container.ID)
			if err == nil {
				fmt.Printf("Successfully killed  %v.\n", Container.Names[0])
			}

			err = model.DockerRemove(Container.ID)
			if err == nil {
				fmt.Printf("Successfully removed %v.\n", Container.Names[0])
			}
		}
	}

	for _, network := range c.Networks {
		if s, _ := model.DockerNetworkStatus(&network); s {
			fmt.Println(s)
			e := model.DockerNetworkRemove(&network)
			if e != nil {
				fmt.Println(e)
			}
		}
		if s, _ := model.DockerNetworkStatus(&network); !s {
			fmt.Printf("Successfully removed network %v\n", network.Name)
		} else {
			fmt.Printf("Network %v was not removed\n", network.Name)
		}
	}

	for _, resolver := range c.Resolvers {
		resolver.Clean()
	}
}
