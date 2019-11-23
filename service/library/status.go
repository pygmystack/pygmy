package library

import (
	"fmt"
	"strings"

	"github.com/fubarhouse/pygmy-go/service/haproxy_connector"
	model "github.com/fubarhouse/pygmy-go/service/interface"
	"github.com/fubarhouse/pygmy-go/service/network"
	"github.com/fubarhouse/pygmy-go/service/resolv"
)

func Status(c Config) {

	Setup(&c)
	checks := DryRun(&c)

	if len(checks) > 0 {
		fmt.Println("Port allocation issue(s) were identified:")
		for _, check := range checks {
			fmt.Println(check.Message)
		}
		fmt.Println()
	}

	// Logic for containers when containers are running.
	Containers, _ := model.DockerContainerList()
	for _, Container := range Containers {
		if Container.Labels["pygmy"] == "pygmy" {
			name := strings.TrimLeft(Container.Names[0], "/")
			for x, Service := range c.Services {
				if Service.Name == name {
					name = x
				}
			}
			Service := c.Services[name]
			if Service.Name != "" {
				if !Service.Disabled && !Service.Discrete && Service.Name != "" {
					if s, _ := Service.Status(); s {
						fmt.Printf("[*] %v: Running as container %v\n", name, Service.Name)
					} else {
						fmt.Printf("[ ] %v is not running\n", Service.Name)
					}
				}
			} else {
				fmt.Printf("[!] %v: Still running as (no longer configured)\n", name)
			}
		}
	}

	// Logic for containers when they're not running.
	for _, Container := range c.Services {
		if s, _ := Container.Status(); !s {
			fmt.Printf("[ ] %v is not running\n", Container.Name)
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

	for _, Container := range c.Services {
		if Container.Group == "showkeys" {
			Container.Output = true
			Container.Start()
		}
	}

}
