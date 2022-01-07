package library

import (
	"fmt"
	"strings"

	. "github.com/logrusorgru/aurora"

	"github.com/pygmystack/pygmy/service/color"
	"github.com/pygmystack/pygmy/service/endpoint"
	model "github.com/pygmystack/pygmy/service/interface"
	"github.com/pygmystack/pygmy/service/interface/docker"
	"github.com/pygmystack/pygmy/service/resolv"
)

// Status will show the state of all the things Pygmy manages.
func Status(c Config) {

	Setup(&c)
	checks := DryRun(&c)
	agentPresent := false

	if len(checks) > 0 {
		for _, check := range checks {
			if check.State {
				color.Print(Green(check.Message + "\n"))
			} else {
				color.Print(Red(check.Message + "\n"))
			}
		}
		fmt.Println()
	}

	Containers, _ := docker.DockerContainerList()
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
							color.Print(Green(fmt.Sprintf("[*] %s: Running as container %s\n", name, name)))
						} else {
							color.Print(Red(fmt.Sprintf("[ ] %s is not running\n", name)))
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
				color.Print(Red(fmt.Sprintf("[ ] %s is not running\n", name)))
			}
		}
	}

	for _, Network := range c.Networks {
		for _, Container := range Network.Containers {
			if x, _ := docker.DockerNetworkConnected(Network.Name, Container.Name); !x {
				color.Print(Red(fmt.Sprintf("[ ] %s is not connected to network %s\n", Container.Name, Network.Name)))
			} else {
				color.Print(Green(fmt.Sprintf("[*] %s is connected to network %s\n", Container.Name, Network.Name)))
			}
		}
	}

	for _, resolver := range c.Resolvers {
		r := resolv.Resolv{Name: resolver.Name, Data: resolver.Data, Folder: resolver.Folder, File: resolver.File}
		if s := r.Status(&model.Params{Domain: c.Domain}); s {
			color.Print(Green(fmt.Sprintf("[*] Resolv %v is properly connected\n", resolver.Name)))
		} else {
			color.Print(Red(fmt.Sprintf("[ ] Resolv %v is not properly connected\n", resolver.Name)))
		}
	}

	for _, volume := range c.Volumes {
		if s, _ := docker.DockerVolumeExists(volume); s {
			color.Print(Green(fmt.Sprintf("[*] Volume %s has been created\n", volume.Name)))
		} else {
			color.Print(Green(fmt.Sprintf("[ ] Volume %s has not ben created\n", volume.Name)))
		}
	}

	// Show ssh-keys in the agent
	if agentPresent {
		for _, v := range c.Services {
			purpose, _ := v.GetFieldString("purpose")
			if purpose == "sshagent" {
				l, _ := docker.DockerExec(v.Config.Labels["pygmy.name"], "ssh-add -l")
				output := strings.Trim(string(l), "\n")
				fmt.Println(output)
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

	// List out all running projects to get their URL.
	containers, _ := docker.DockerContainerList()
	var urls []string
	for _, container := range containers {
		if container.State == "running" && !strings.Contains(fmt.Sprint(container.Names), "amazeeio") {
			obj, _ := docker.DockerInspect(container.ID)
			vars := obj.Config.Env
			for _, v := range vars {
				// Look for the environment variable $LAGOON_ROUTE.
				if strings.Contains(v, "LAGOON_ROUTE=") {
					url := strings.TrimPrefix(v, "LAGOON_ROUTE=")
					if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
						url = "http://" + url
					}
					urls = append(urls, url)
				}
			}
		}
	}

	cleanurls := unique(urls)
	for _, url := range cleanurls {
		endpoint.Validate(url)
		if r := endpoint.Validate(url); r {
			fmt.Printf(" - %v\n", url)
		} else {
			fmt.Printf(" ! %v\n", url)
		}
	}

}
