package library

import (
	"fmt"

	"github.com/fubarhouse/pygmy/v1/service/haproxy_connector"
	model "github.com/fubarhouse/pygmy/v1/service/interface"
	"github.com/fubarhouse/pygmy/v1/service/network"
	"github.com/fubarhouse/pygmy/v1/service/resolv"
)

func Status(c Config) {

	Setup(&c)

	for Label, Service := range c.Services {
		if !Service.Disabled && !Service.Discrete {
			if s, _ := Service.Status(); s {
				fmt.Printf("[*] %v: Running as container %v\n", Label, Service.Name)
			} else {
				fmt.Printf("[ ] %v is not running\n", Label)
			}
		}
	}

	{
		Containers, _ := model.DockerContainerList()
		for _, Container := range Containers {
			if Container.State == "running" {
				if e := c.Services[fmt.Sprint(Container.Names)]; e.Name == "" {
					if f := Container.Labels["pygmy"]; f == "pygmy" {
						fmt.Printf("[!] %v: Still running as container %v\n", Container.Image, Container.Names[0])
					}
				}
			}
		}
	}

	for Network, Containers := range c.Networks {
		netStat, _ := network.Status(Network)
		if netStat {
			for _, Container := range Containers {
				if s, _ := haproxy_connector.Connected(Container, Network); s {
					fmt.Printf("[*] %v is connected to network %v\n", Container, Network)
				} else {
					fmt.Printf("[ ] %v is not connected to network %v\n", Container, Network)
				}
			}
		}
	}

	for _, resolver := range c.Resolvers {
		r := resolv.New(resolv.Resolv{Name: resolver.Name, Data: resolver.Data, Folder: resolver.Folder, File: resolver.File})
		if s := r.Status(); s {
			fmt.Printf("[*] Resolv %v is properly connected\n", resolver.Name)
		} else {
			fmt.Printf("[ ] Resolv %v is not properly connected\n", resolver.Name)
		}
	}

	for _, volume := range c.Volumes {
		if s, _ := model.DockerVolumeExists(volume); s {
			fmt.Printf("[*] Volume %v has been created\n", volume)
		} else {
			fmt.Printf("[ ] Volume %v has not been created\n", volume)
		}
	}

}
