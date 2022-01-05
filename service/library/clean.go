package library

import (
	"fmt"
	"strings"

	. "github.com/logrusorgru/aurora"

	"github.com/pygmystack/pygmy/service/color"
	"github.com/pygmystack/pygmy/service/interface/docker"
)

// Clean will forcibly kill and remove all of pygmy's containers in the daemon
func Clean(c Config) {

	Setup(&c)
	Containers, _ := docker.DockerContainerList()
	NetworksToClean := []string{}

	for _, Container := range Containers {
		ContainerName := strings.Trim(Container.Names[0], "/")
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
				color.Print(Green(fmt.Sprintf("Successfully killed %s\n", ContainerName)))
			}

			err = docker.DockerRemove(Container.ID)
			if err == nil {
				color.Print(Green(fmt.Sprintf("Successfully removed %s\n", ContainerName)))
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
				color.Print(Green(fmt.Sprintf("Successfully removed network %s\n", NetworksToClean[n])))
			} else {
				color.Print(Red(fmt.Sprintf("Failed to remove %s\n", NetworksToClean[n])))
			}
		}
	}

	for _, resolver := range c.Resolvers {
		resolver.Clean()
	}
}
