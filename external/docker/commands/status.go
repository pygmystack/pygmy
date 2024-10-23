package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"

	"github.com/pygmystack/pygmy/internal/runtime"
	runtimecontainers "github.com/pygmystack/pygmy/internal/runtime/docker/internals/containers"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/networks"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/volumes"
	"github.com/pygmystack/pygmy/internal/utils/color"
	"github.com/pygmystack/pygmy/internal/utils/endpoint"
	"github.com/pygmystack/pygmy/internal/utils/resolv"
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

	// Ensure the services struct is not nil.
	c.JSONStatus.Services = make(map[string]StatusJSONStatus)

	Containers, _ := runtimecontainers.List()
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
							c.JSONStatus.Services[name] = StatusJSONStatus{
								Container: name,
								ImageRef:  Service.Image,
								State:     true,
							}
						} else {
							c.JSONStatus.Services[name] = StatusJSONStatus{
								Container: name,
								ImageRef:  Service.Image,
							}
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
				c.JSONStatus.Services[name] = StatusJSONStatus{
					Container: name,
					ImageRef:  Service.Image,
				}
			}
		}
	}

	for _, Network := range c.Networks {
		for _, Container := range Network.Containers {
			if x, _ := networks.Connected(Network.Name, Container.Name); !x {
				c.JSONStatus.Networks = append(c.JSONStatus.Networks, fmt.Sprintf("%s is not connected to the network %s", Container.Name, Network.Name))
			} else {
				c.JSONStatus.Networks = append(c.JSONStatus.Networks, fmt.Sprintf("%s is connected to the network %s", Container.Name, Network.Name))
			}
		}
	}

	for _, resolver := range c.Resolvers {
		r := resolv.Resolv{Name: resolver.Name, Data: resolver.Data, Folder: resolver.Folder, File: resolver.File}
		if s := r.Status(&runtime.Params{Domain: c.Domain}); s {
			c.JSONStatus.Resolvers = append(c.JSONStatus.Resolvers, fmt.Sprintf("Resolv %s is properly connected", resolver.Name))
		} else {
			c.JSONStatus.Resolvers = append(c.JSONStatus.Resolvers, fmt.Sprintf("Resolv %s is not properly connected", resolver.Name))
		}
	}

	for _, volume := range c.Volumes {
		if s, _ := volumes.Exists(volume.Name); s {
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
				l, _ := runtimecontainers.Exec(v.Config.Labels["pygmy.name"], "ssh-add -l")
				// Remove \u0000 & \u0001 from output messages.
				output := strings.ReplaceAll(string(l), "\u0000", "")
				output = strings.ReplaceAll(output, "\u0001", "")
				output = strings.Trim(output, "\n")
				c.JSONStatus.SSHMessages = append(c.JSONStatus.SSHMessages, output)
			}
		}
	}

	// List out all running projects to get their URL.
	var urls []string

	for _, Container := range c.Services {
		Status, _ := Container.Status()
		url, _ := Container.GetFieldString("url")
		if url != "" && Status {
			urls = append(urls, url)
		}
	}

	containers, _ := runtimecontainers.List()
	for _, container := range containers {
		if container.State == "running" && !strings.Contains(fmt.Sprint(container.Names), "amazeeio") {
			obj, _ := runtimecontainers.Inspect(container.ID)
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

	for k, v := range c.JSONStatus.Services {
		if v.State {
			color.Print(aurora.Green(fmt.Sprintf("[*] %s: Running as container %s\n", k, v.Container)))
		} else {
			color.Print(aurora.Red(fmt.Sprintf("[ ] %s is not running\n", k)))
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
