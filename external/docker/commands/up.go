package commands

import (
	"fmt"
	"os"
	"strings"

	. "github.com/logrusorgru/aurora"

	"github.com/pygmystack/pygmy/external/docker/setup"
	"github.com/pygmystack/pygmy/internal/runtime/docker"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
	runtimecontainers "github.com/pygmystack/pygmy/internal/runtime/docker/internals/containers"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/networks"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/volumes"
	"github.com/pygmystack/pygmy/internal/utils/color"
	"github.com/pygmystack/pygmy/internal/utils/endpoint"
)

// Up will bring Pygmy up.
func Up(c setup.Config) error {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return err
	}

	setup.Setup(ctx, cli, &c)
	checks, _ := setup.DryRun(ctx, cli, &c)
	agentPresent := false

	foundIssues := 0
	for _, check := range checks {
		if !check.State {
			fmt.Println(check.Message)
			foundIssues++
		}
	}
	if foundIssues > 0 {
		fmt.Println("Please address the above issues before you attempt to start Pygmy again.")
		os.Exit(1)
	}

	for _, volume := range c.Volumes {
		if s, _ := volumes.Exists(ctx, cli, volume.Name); !s {
			_, err := volumes.Create(ctx, cli, volume)
			if err == nil {
				color.Print(Green(fmt.Sprintf("Created volume %s\n", volume.Name)))
			} else {
				fmt.Println(err)
			}
		} else {
			color.Print(Green(fmt.Sprintf("Already created volume %s\n", volume.Name)))
		}
	}

	// Maps are... bad for predictable sequencing.
	// Look over the sorted slice and start them in
	// alphabetical order - so that one can configure
	// an ssh-agent like amazeeio-ssh-agent.
	for _, s := range c.SortedServices {
		service := c.Services[s]
		enabled, _ := service.GetFieldBool(ctx, cli, "enable")
		purpose, _ := service.GetFieldString(ctx, cli, "purpose")
		name, _ := service.GetFieldString(ctx, cli, "name")

		// Do not show or add keys:
		if enabled && purpose != "addkeys" {

			if se := service.Setup(ctx, cli); se == nil {
				fmt.Print(Green(fmt.Sprintf("Successfully pulled %s\n", service.Config.Image)))
			}
			if status, _ := service.Status(ctx, cli); !status {
				if ce := service.Create(ctx, cli); ce != nil {
					// If the status is false but the container is already created, we can ignore that error.
					if !strings.Contains(ce.Error(), "namespace is already taken") {
						fmt.Printf("Failed to create %s: %s\n", Red(name), ce)
					}
				}
				if se := service.Start(ctx, cli); se == nil {
					fmt.Print(Green(fmt.Sprintf("Successfully started %s\n", name)))
				} else {
					fmt.Printf("Failed to start %s: %s\n", Red(name), se)
				}
			} else {
				fmt.Print(Green(fmt.Sprintf("Already started %s\n", name)))
			}
		}

		// If one or more agent was found:
		if purpose == "sshagent" {
			agentPresent = true
		}
	}

	// Docker network(s) creation
	for _, Network := range c.Networks {
		if Network.Name != "" {
			netVal, _ := networks.Status(ctx, cli, Network.Name)
			if !netVal {
				if err := networks.Create(ctx, cli, &Network); err == nil {
					color.Print(Green(fmt.Sprintf("Successfully created network %s\n", Network.Name)))
				} else {
					color.Print(Red(fmt.Sprintf("Could not create network %s\n", Network.Name)))
				}
			}
		}
	}

	// Container network connection(s)
	for _, s := range c.SortedServices {
		service := c.Services[s]
		name, nameErr := service.GetFieldString(ctx, cli, "name")
		// If the network is configured at the container level, connect it.
		if Network, _ := service.GetFieldString(ctx, cli, "network"); Network != "" && nameErr == nil {
			if s, _ := networks.Connected(ctx, cli, Network, name); !s {
				if s := networks.Connect(ctx, cli, Network, name); s == nil {
					color.Print(Green(fmt.Sprintf("Successfully connected %s to %s\n", name, Network)))
				} else {
					discrete, _ := service.GetFieldBool(ctx, cli, "discrete")
					if !discrete {
						color.Print(Red(fmt.Sprintf("Could not connect %s to %s\n", name, Network)))
					}
				}
			} else {
				color.Print(Green(fmt.Sprintf("Already connected %s to %s\n", name, Network)))
			}
		}
	}

	for _, resolver := range c.Resolvers {
		if !c.ResolversDisabled {
			if !resolver.Status(&docker.Params{Domain: c.Domain}) {
				resolver.Configure(&docker.Params{Domain: c.Domain})
			}
		}
	}

	// Add ssh-keys to the agent
	if agentPresent {
		for _, v := range c.Keys {
			if e := SshKeyAdd(c, v.Path); e != nil {
				color.Print(Red(fmt.Sprintf("%v\n", e)))
			}
		}
	}

	for _, service := range c.Services {
		name, _ := service.GetFieldString(ctx, cli, "name")
		url, _ := service.GetFieldString(ctx, cli, "url")
		if s, _ := service.Status(ctx, cli); s && url != "" {
			endpoint.Validate(url)
			if r := endpoint.Validate(url); r {
				fmt.Printf(" - %v (%v)\n", url, name)
			} else {
				fmt.Printf(" ! %v (%v)\n", url, name)
			}
		}
	}

	// List out all running projects to get their URL.
	containers, _ := runtimecontainers.List(ctx, cli)
	var urls []string
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
		endpoint.Validate(url)
		if r := endpoint.Validate(url); r {
			fmt.Printf(" - %v\n", url)
		} else {
			fmt.Printf(" ! %v\n", url)
		}
	}

	return nil
}
