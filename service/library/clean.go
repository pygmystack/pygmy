package library

import (
	"fmt"

	"github.com/fubarhouse/pygmy-go/service/interface/docker"
	. "github.com/logrusorgru/aurora"
)

// Clean will forcibly kill and remove all of pygmy's containers in the daemon
func Clean(c Config) {

	Setup(&c)
	Containers, _ := docker.DockerContainerList()
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
			err := docker.DockerKill(Container.ID)
			if err == nil {
				fmt.Print(Green(fmt.Sprintf("Successfully killed %s\n", Container.Names[0])))
			}

			err = docker.DockerRemove(Container.ID)
			if err == nil {
				fmt.Print(Green(fmt.Sprintf("Successfully removed %s\n", Container.Names[0])))
			}
		}
	}

	for _, network := range c.Networks {
		NetworksToClean = append(NetworksToClean, network.Name)
	}

	for n := range unique(NetworksToClean) {
		if s, _ := docker.DockerNetworkStatus(NetworksToClean[n]); s {
			e := docker.DockerNetworkRemove(NetworksToClean[n])
			if e != nil {
				fmt.Println(e)
			}
			if s, _ := docker.DockerNetworkStatus(NetworksToClean[n]); !s {
				fmt.Print(Green(fmt.Sprintf("Successfully removed network %s\n", NetworksToClean[n])))
			} else {
				fmt.Print(Red(fmt.Sprintf("Successfully started %s\n", NetworksToClean[n])))
			}
		}
	}

	for _, resolver := range c.Resolvers {
		resolver.Clean()
	}
}
