package commands

import (
	"fmt"
	"strings"

	. "github.com/logrusorgru/aurora"

	"github.com/pygmystack/pygmy/internal/runtimes/docker/containers"
	"github.com/pygmystack/pygmy/internal/runtimes/docker/networks"
	"github.com/pygmystack/pygmy/service/color"
)

// Clean will forcibly kill and remove all of pygmy's containers in the daemon
func Clean(c Config) {

	Setup(&c)
	Containers, _ := containers.List()
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
			err := containers.Kill(Container.ID)
			if err == nil {
				color.Print(Green(fmt.Sprintf("Successfully killed %s\n", ContainerName)))
			}

			err = containers.Remove(Container.ID)
			if err == nil {
				color.Print(Green(fmt.Sprintf("Successfully removed %s\n", ContainerName)))
			}
		}
	}

	for _, network := range c.Networks {
		NetworksToClean = append(NetworksToClean, network.Name)
	}

	for n := range unique(NetworksToClean) {
		if s, _ := networks.Status(NetworksToClean[n]); s {
			e := networks.Remove(NetworksToClean[n])
			if e != nil {
				fmt.Println(e)
			}
			if s, _ := networks.Status(NetworksToClean[n]); !s {
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
