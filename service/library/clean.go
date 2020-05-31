package library

import (
	"fmt"

	model "github.com/fubarhouse/pygmy-go/service/interface"
)

// Clean will forcibly kill and remove all of pygmy's containers in the daemon
func Clean(c Config) {

	Setup(&c)
	Containers, _ := model.DockerContainerList()
	NetworksToClean := []string{}

	for _, Container := range Containers {
		target := false
		if l := Container.Labels["pygmy.enable"]; l == "true" || l == "1" {
			target = true
		}
		if l := Container.Labels["pygmy"]; l == "pygmy" {
			target = true
		}
		if l := Container.Labels["pygmy.network"]; l != "" {
			NetworksToClean = append(NetworksToClean, l)
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
		NetworksToClean = append(NetworksToClean, network.Name)
	}

	for n := range unique(NetworksToClean) {
		if s, _ := model.DockerNetworkStatus(NetworksToClean[n]); s {
			e := model.DockerNetworkRemove(NetworksToClean[n])
			if e != nil {
				fmt.Println(e)
			}
			if s, _ := model.DockerNetworkStatus(NetworksToClean[n]); !s {
				fmt.Printf("Successfully removed network %v\n", NetworksToClean[n])
			} else {
				fmt.Printf("Network %v was not removed\n", NetworksToClean[n])
			}
		}
	}

	for _, resolver := range c.Resolvers {
		resolver.Clean()
	}
}
