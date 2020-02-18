package library

import (
	"fmt"
	"strings"

	"github.com/fubarhouse/pygmy-go/service/endpoint"
	model "github.com/fubarhouse/pygmy-go/service/interface"
	"github.com/fubarhouse/pygmy-go/service/resolv"
)

// Status will show the state of all the things Pygmy manages.
func Status(c Config) {

	Setup(&c)
	checks := DryRun(&c)
	agentPresent := false

	if len(checks) > 0 {
		for _, check := range checks {
			fmt.Println(check.Message)
		}
		fmt.Println()
	}

	Containers, _ := model.DockerContainerList()
	for _, Container := range Containers {
		if Container.Labels["pygmy.enable"] == "true" || Container.Labels["pygmy.enable"] == "1" {
			Service := c.Services[strings.Trim(Container.Names[0], "/")]
			if s, _ := Service.Status(); s {
				name, _ := Service.GetFieldString("name")
				enabled, _ := Service.GetFieldBool("enable")
				discrete, _ := Service.GetFieldBool("discrete")
				purpose, _ := Service.GetFieldString("purpose")
				if name != "" {
					if purpose == "sshagent" {
						agentPresent = true
					}
					if enabled && !discrete && name != "" {
						if s, _ := Service.Status(); s {
							fmt.Printf("[*] %v: Running as container %v\n", name, name)
						} else {
							fmt.Printf("[ ] %v is not running\n", name)
						}
					}
				}
			}
		}
	}

	for _, Service := range c.Services {
		if s, _ := Service.Status(); !s {
			name, _ := Service.GetFieldString("name")
			discrete, _ := Service.GetFieldBool("discrete")
			if !discrete {
				fmt.Printf("[ ] %v is not running\n", name)
			}
		}
	}

	for _, Network := range c.Networks {
		for _, Container := range Network.Containers {
			if x, _ := model.DockerNetworkConnected(Network, Container.Name); !x {
				fmt.Printf("[ ] %v is not connected to network %v\n", Container.Name, Network.Name)
			} else {
				fmt.Printf("[*] %v is connected to network %v\n", Container.Name, Network.Name)
			}
		}
	}

	for _, resolver := range c.Resolvers {
		r := resolv.Resolv{Name: resolver.Name, Data: resolver.Data, Folder: resolver.Folder, File: resolver.File}
		if s := r.Status(); s {
			fmt.Printf("[*] Resolv %v is properly connected\n", resolver.Name)
		} else {
			fmt.Printf("[ ] Resolv %v is not properly connected\n", resolver.Name)
		}
	}

	for _, volume := range c.Volumes {
		if s, _ := model.DockerVolumeExists(volume); s {
			fmt.Printf("[*] Volume %v has been created\n", volume.Name)
		} else {
			fmt.Printf("[ ] Volume %v has not been created\n", volume.Name)
		}
	}

	// Show ssh-keys in the agent
	if agentPresent {
		for _, v := range c.Services {
			purpose, _ := v.GetFieldString("purpose")
			if purpose == "showkeys" {
				out, _ := v.Start()
				if len(string(out)) > 0 {
					output := strings.Trim(string(out), "\n")
					fmt.Println(output)
				}
			}
		}
	}

	for _, Container := range c.Services {
		Status, _ := Container.Status()
		name, _ := Container.GetFieldString("name")
		url, _ := Container.GetFieldString("url")
		if url != "" && Status {
			endpoint.Validate(url)
			if r := endpoint.Validate(url); r {
				fmt.Printf(" - %v (%v)\n", url, name)
			} else {
				fmt.Printf(" ! %v (%v)\n", url, name)
			}
		}
	}

}
