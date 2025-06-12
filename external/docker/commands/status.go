package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	urltools "net/url"
	"os"
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
						start(c.Debug, fmt.Sprintf("Checking Service for SSH Role: %s", Container.Names[0]))
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
						finish(c.Debug, fmt.Sprintf("Checking Service for SSH Role: %s", Container.Names[0]))
					}
				}
			}
		}
	}

	for n, Service := range c.Services {
		start(c.Debug, fmt.Sprintf("Checking Service Status: %s", n))
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
		finish(c.Debug, fmt.Sprintf("Checking Service Status: %s", n))
	}

	for _, Network := range c.Networks {
		start(c.Debug, fmt.Sprintf("Checking Network Connection: %s", Network.Name))
		for _, Container := range Network.Containers {
			start(c.Debug, fmt.Sprintf("Checking Network Connection Status: %s", Container.Name))
			if x, _ := networks.Connected(ctx, cli, Network.Name, Container.Name); !x {
				c.JSONStatus.Networks = append(c.JSONStatus.Networks, fmt.Sprintf("%s is not connected to the network %s", Container.Name, Network.Name))
			} else {
				c.JSONStatus.Networks = append(c.JSONStatus.Networks, fmt.Sprintf("%s is connected to the network %s", Container.Name, Network.Name))
			}
			finish(c.Debug, fmt.Sprintf("Checking Network Connection Status: %s", Container.Name))
		}
		finish(c.Debug, fmt.Sprintf("Checking Network Connection: %s", Network.Name))
	}

	for _, resolver := range c.Resolvers {
		start(c.Debug, fmt.Sprintf("Checking Resolver Status: %s", resolver.Name))
		r := resolv.Resolv{Name: resolver.Name, Data: resolver.Data, Folder: resolver.Folder, File: resolver.File}
		if s := r.Status(&docker.Params{Domain: c.Domain}); s {
			c.JSONStatus.Resolvers = append(c.JSONStatus.Resolvers, fmt.Sprintf("Resolv %s is properly connected", resolver.Name))
		} else {
			c.JSONStatus.Resolvers = append(c.JSONStatus.Resolvers, fmt.Sprintf("Resolv %s is not properly connected", resolver.Name))
		}
		finish(c.Debug, fmt.Sprintf("Checking Resolver Status: %s", resolver.Name))
	}

	for _, volume := range c.Volumes {
		start(c.Debug, fmt.Sprintf("Checking Volume Status: %s", volume.Name))
		if s, _ := volumes.Exists(ctx, cli, volume.Name); s {
			c.JSONStatus.Volumes = append(c.JSONStatus.Volumes, fmt.Sprintf("Volume %s has been created", volume.Name))
		} else {
			c.JSONStatus.Volumes = append(c.JSONStatus.Volumes, fmt.Sprintf("Volume %s has not been created", volume.Name))
		}
		finish(c.Debug, fmt.Sprintf("Checking Volume Status: %s", volume.Name))
	}

	// Show ssh-keys in the agent
	if agentPresent {
		for _, v := range c.Services {
			start(c.Debug, fmt.Sprintf("Checking for SSH Agent: %s", v.Config.Labels["pygmy.name"]))
			purpose, _ := v.GetFieldString(ctx, cli, "purpose")
			if purpose == "sshagent" {
				l, _ := runtimecontainers.Exec(ctx, cli, v.Config.Labels["pygmy.name"], "ssh-add -l")
				// Remove \u0000 & \u0001 from output messages.
				output := strings.ReplaceAll(string(l), "\u0000", "")
				output = strings.ReplaceAll(output, "\u0001", "")
				output = strings.Trim(output, "\n")
				c.JSONStatus.SSHMessages = append(c.JSONStatus.SSHMessages, output)
			}
			finish(c.Debug, fmt.Sprintf("Checking for SSH Agent: %s", v.Config.Labels["pygmy.name"]))
		}
	}

	// List out all running projects to get their URL.
	var urls []string

	for _, Container := range c.Services {
		Status, _ := Container.Status(ctx, cli)
		url, _ := Container.GetFieldString(ctx, cli, "url")

		cleanUrl, err := urltools.Parse(url)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if Status {
			if strings.Contains(cleanUrl.String(), "docker.amazee.io") {
				finalUrl := strings.Trim(cleanUrl.String(), "[]")
				fmt.Printf("Added %s\n", finalUrl)
				urls = append(urls, finalUrl)
			}
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
		} else {
			finish(c.Debug, fmt.Sprintf("Not a Pygmy Container: %s", container.ID))
		}
		finish(c.Debug, fmt.Sprintf("Inspecting Container Status: %s", container.ID))
	}

	cleanurls := setup.Unique(urls)
	for _, url := range cleanurls {
		start(c.Debug, fmt.Sprintf("Validating Endpoint: %s", url))

		cleanUrl, err := urltools.Parse(url)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if os.Getenv("PYGMY_WEB_PORT") != "" {
			cleanUrl.Host = net.JoinHostPort(cleanUrl.Host, os.Getenv("PYGMY_WEB_PORT"))
		}
		finalUrl := strings.Trim(cleanUrl.String(), "[]")

		result, statuscode := endpoint.Validate(finalUrl)
		c.JSONStatus.URLValidations = append(c.JSONStatus.URLValidations, setup.StatusJSONURLValidation{
			Endpoint:   finalUrl,
			Success:    result,
			StatusCode: statuscode,
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
			fmt.Printf(" ! %s (Status Code %v)\n", v.Endpoint, v.StatusCode)
		}
	}

}
