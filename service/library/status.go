package library

import (
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora"
	"strings"

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
			c.JSONStatus.PortAvailability = append(c.JSONStatus.PortAvailability, check.Message)
		}
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
						st := StatusJSONStatus{
							Name:      name,
							Container: name,
						}
						if s, _ := Service.Status(); s {
							st.State = "running"
							c.JSONStatus.Services = append(c.JSONStatus.Services, st)
						} else {
							st.State = "not running"
							c.JSONStatus.Services = append(c.JSONStatus.Services, st)
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
				st := StatusJSONStatus{
					Name:      name,
					Container: name,
					State:     "not running",
				}
				c.JSONStatus.Services = append(c.JSONStatus.Services, st)
			}
		}
	}

	for _, Network := range c.Networks {
		for _, Container := range Network.Containers {
			if x, _ := docker.DockerNetworkConnected(Network.Name, Container.Name); !x {
				c.JSONStatus.Networks = append(c.JSONStatus.Networks, fmt.Sprintf("%s is not connected to the network %s", Container.Name, Network.Name))
			} else {
				c.JSONStatus.Networks = append(c.JSONStatus.Networks, fmt.Sprintf("%s is connected to the network %s", Container.Name, Network.Name))
			}
		}
	}

	for _, resolver := range c.Resolvers {
		r := resolv.Resolv{Name: resolver.Name, Data: resolver.Data, Folder: resolver.Folder, File: resolver.File}
		if s := r.Status(&model.Params{Domain: c.Domain}); s {
			c.JSONStatus.Resolvers = append(c.JSONStatus.Resolvers, fmt.Sprintf("Resolv %s is properly connected", resolver.Name))
		} else {
			c.JSONStatus.Resolvers = append(c.JSONStatus.Resolvers, fmt.Sprintf("Resolv %s is not properly connected", resolver.Name))
		}
	}

	for _, volume := range c.Volumes {
		if s, _ := docker.DockerVolumeExists(volume.Name); s {
			c.JSONStatus.Volumes = append(c.JSONStatus.Volumes, fmt.Sprintf("Volume %s has been created", volume.Name))
		} else {
			c.JSONStatus.Volumes = append(c.JSONStatus.Volumes, fmt.Sprintf("Volume %s has not been created", volume.Name))
		}
	}

	// Show ssh-keys in the agent
	if agentPresent {
		for _, v := range c.Services {
			purpose, _ := v.GetFieldString("purpose")
			if purpose == "sshagent" {
				l, _ := docker.DockerExec(v.Config.Labels["pygmy.name"], "ssh-add -l")
				// Remove \u0000 & \u0001 from output messages.
				output := strings.ReplaceAll(string(l), "\u0000", "")
				output = strings.ReplaceAll(output, "\u0001", "")
				output = strings.Trim(output, "\n")
				c.JSONStatus.SSHMessages = append(c.JSONStatus.SSHMessages, output)
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
		if r := endpoint.Validate(url); !r {
			c.JSONStatus.URLValidations = append(c.JSONStatus.URLValidations, fmt.Sprintf(" ! %v\n", url))
		} else {
			c.JSONStatus.URLValidations = append(c.JSONStatus.URLValidations, fmt.Sprintf(" - %v\n", url))
		}
	}

	if c.JSONFormat {
		PrintStatusJSON(c)
		return
	}

	PrintStatusHumanReadable(c)

}

func PrintStatusJSON(c Config) {
	jsonData, _ := json.Marshal(c.JSONStatus)
	fmt.Println(string(jsonData))

}
func PrintStatusHumanReadable(c Config) {
	for _, v := range c.JSONStatus.PortAvailability {
		if strings.Contains(v, "is not able to start on port") {
			color.Print(aurora.Red(fmt.Sprintf("[ ] %s\n", v)))
		} else {
			color.Print(aurora.Green(fmt.Sprintf("[*] %s\n", v)))
		}
	}

	for _, v := range c.JSONStatus.Services {
		if strings.Contains(v.State, "not running") {
			color.Print(aurora.Red(fmt.Sprintf("[ ] %s is not running\n", v.Name)))
		} else {
			color.Print(aurora.Green(fmt.Sprintf("[*] %s: Running as container %s\n", v.Name, v.Container)))
		}
	}

	for _, v := range c.JSONStatus.Resolvers {
		if strings.Contains(v, "not properly connected") {
			color.Print(aurora.Red(fmt.Sprintf("[ ] %s\n", v)))
		} else {
			color.Print(aurora.Green(fmt.Sprintf("[*] %s\n", v)))
		}
	}

	for _, v := range c.JSONStatus.Networks {
		if strings.Contains(v, "is not connected to network") {
			color.Print(aurora.Red(fmt.Sprintf("[ ] %s\n", v)))
		} else {
			color.Print(aurora.Green(fmt.Sprintf("[*] %s\n", v)))
		}
	}

	for _, v := range c.JSONStatus.Volumes {
		if strings.Contains(v, "has not been created") {
			color.Print(aurora.Red(fmt.Sprintf("[ ] %s\n", v)))
		} else {
			color.Print(aurora.Green(fmt.Sprintf("[*] %s\n", v)))
		}
	}

	for _, v := range c.JSONStatus.SSHMessages {
		fmt.Printf("%s\n", v)
	}

	for _, v := range c.JSONStatus.URLValidations {
		fmt.Printf("%s\n", v)
	}

}
