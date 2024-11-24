package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	"github.com/logrusorgru/aurora"

	"github.com/pygmystack/pygmy/external/docker/setup"
	"github.com/pygmystack/pygmy/internal/runtime/docker"
	runtimecontainers "github.com/pygmystack/pygmy/internal/runtime/docker/internals/containers"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/networks"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/volumes"
	"github.com/pygmystack/pygmy/internal/utils/color"
	"github.com/pygmystack/pygmy/internal/utils/endpoint"
	"github.com/pygmystack/pygmy/internal/utils/resolv"
)

// Status will show the state of all the things Pygmy manages.
func Status(ctx context.Context, cli *client.Client, c setup.Config) {
	setup.Setup(ctx, cli, &c)
	checks, _ := setup.DryRun(ctx, cli, &c)
	agentPresent := false

	if len(checks) > 0 {
		for _, check := range checks {
			c.JSONStatus.PortAvailability = append(c.JSONStatus.PortAvailability, check.Message)
		}
	}

	// Ensure the services struct is not nil.
	c.JSONStatus.Services = make(map[string]setup.StatusJSONStatus)

	Containers, _ := runtimecontainers.List(ctx, cli)
	for _, Container := range Containers {
		if Container.Labels["pygmy.enable"] == "true" || Container.Labels["pygmy.enable"] == "1" {
			Service := c.Services[strings.Trim(Container.Names[0], "/")]
			if s, _ := Service.Status(ctx, cli); s {
				name, _ := Service.GetFieldString(ctx, cli, "name")
				enabled, _ := Service.GetFieldBool(ctx, cli, "enable")
				discrete, _ := Service.GetFieldBool(ctx, cli, "discrete")
				purpose, _ := Service.GetFieldString(ctx, cli, "purpose")
				if name != "" {
					if purpose == "sshagent" {
						agentPresent = true
					}
					if enabled && !discrete && name != "" {
						if s, _ := Service.Status(ctx, cli); s {
							c.JSONStatus.Services[name] = setup.StatusJSONStatus{
								Container: name,
								ImageRef:  Service.Image,
								State:     true,
							}
						} else {
							c.JSONStatus.Services[name] = setup.StatusJSONStatus{
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
		if s, _ := Service.Status(ctx, cli); !s {
			name, _ := Service.GetFieldString(ctx, cli, "name")
			discrete, _ := Service.GetFieldBool(ctx, cli, "discrete")
			if !discrete {
				c.JSONStatus.Services[name] = setup.StatusJSONStatus{
					Container: name,
					ImageRef:  Service.Image,
				}
			}
		}
	}

	for _, Network := range c.Networks {
		for _, Container := range Network.Containers {
			if x, _ := networks.Connected(ctx, cli, Network.Name, Container.Name); !x {
				c.JSONStatus.Networks = append(c.JSONStatus.Networks, fmt.Sprintf("%s is not connected to the network %s", Container.Name, Network.Name))
			} else {
				c.JSONStatus.Networks = append(c.JSONStatus.Networks, fmt.Sprintf("%s is connected to the network %s", Container.Name, Network.Name))
			}
		}
	}

	for _, resolver := range c.Resolvers {
		r := resolv.Resolv{Name: resolver.Name, Data: resolver.Data, Folder: resolver.Folder, File: resolver.File}
		if s := r.Status(&docker.Params{Domain: c.Domain}); s {
			c.JSONStatus.Resolvers = append(c.JSONStatus.Resolvers, fmt.Sprintf("Resolv %s is properly connected", resolver.Name))
		} else {
			c.JSONStatus.Resolvers = append(c.JSONStatus.Resolvers, fmt.Sprintf("Resolv %s is not properly connected", resolver.Name))
		}
	}

	for _, volume := range c.Volumes {
		if s, _ := volumes.Exists(ctx, cli, volume.Name); s {
			c.JSONStatus.Volumes = append(c.JSONStatus.Volumes, fmt.Sprintf("Volume %s has been created", volume.Name))
		} else {
			c.JSONStatus.Volumes = append(c.JSONStatus.Volumes, fmt.Sprintf("Volume %s has not been created", volume.Name))
		}
	}

	// Show ssh-keys in the agent
	if agentPresent {
		for _, v := range c.Services {
			purpose, _ := v.GetFieldString(ctx, cli, "purpose")
			if purpose == "sshagent" {
				l, _ := runtimecontainers.Exec(ctx, cli, v.Config.Labels["pygmy.name"], "ssh-add -l")
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
		Status, _ := Container.Status(ctx, cli)
		url, _ := Container.GetFieldString(ctx, cli, "url")
		if url != "" && Status {
			urls = append(urls, url)
		}
	}

	containers, _ := runtimecontainers.List(ctx, cli)
	for _, container := range containers {
		if container.State == "running" && !strings.Contains(fmt.Sprint(container.Names), "amazeeio") {
			obj, _ := runtimecontainers.Inspect(ctx, cli, container.ID)
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

	cleanurls := setup.Unique(urls)
	for _, url := range cleanurls {
		result := endpoint.Validate(url)
		c.JSONStatus.URLValidations = append(c.JSONStatus.URLValidations, setup.StatusJSONURLValidation{
			Endpoint: url,
			Success:  result,
		})
	}

	if c.JSONFormat {
		PrintStatusJSON(c)
		return
	}

	PrintStatusHumanReadable(c)

}

func PrintStatusJSON(c setup.Config) {
	jsonData, _ := json.Marshal(c.JSONStatus)
	fmt.Println(string(jsonData))

}
func PrintStatusHumanReadable(c setup.Config) {
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
		if v.Success {
			fmt.Printf(" - %s\n", v.Endpoint)
		} else {
			fmt.Printf(" ! %s\n", v.Endpoint)
		}
	}

}
