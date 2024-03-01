package library

import (
	"fmt"
	"os"
	"strings"

	. "github.com/logrusorgru/aurora"

	"github.com/pygmystack/pygmy/service/color"
	"github.com/pygmystack/pygmy/service/endpoint"
	model "github.com/pygmystack/pygmy/service/interface"
	"github.com/pygmystack/pygmy/service/interface/docker"
)

// Up will bring Pygmy up.
func Up(c Config) {

	Setup(&c)
	checks := DryRun(&c)
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
		if s, _ := docker.DockerVolumeExists(volume.Name); !s {
			_, err := docker.DockerVolumeCreate(volume)
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
		enabled, _ := service.GetFieldBool("enable")
		purpose, _ := service.GetFieldString("purpose")
		name, _ := service.GetFieldString("name")

		// Do not show or add keys:
		if enabled && purpose != "addkeys" {

			if se := service.Setup(); se == nil {
				fmt.Print(Green(fmt.Sprintf("Successfully pulled %s\n", service.Config.Image)))
			}
			if status, _ := service.Status(); !status {
				if ce := service.Create(); ce != nil {
					// If the status is false but the container is already created, we can ignore that error.
					if !strings.Contains(ce.Error(), "namespace is already taken") {
						fmt.Printf("Failed to create %s: %s\n", Red(name), ce)
					}
				}
				if se := service.Start(); se == nil {
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
			netVal, _ := docker.DockerNetworkStatus(Network.Name)
			if !netVal {
				if err := NetworkCreate(Network); err == nil {
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
		name, nameErr := service.GetFieldString("name")
		// If the network is configured at the container level, connect it.
		if Network, _ := service.GetFieldString("network"); Network != "" && nameErr == nil {
			if s, _ := docker.DockerNetworkConnected(Network, name); !s {
				if s := NetworkConnect(Network, name); s == nil {
					color.Print(Green(fmt.Sprintf("Successfully connected %s to %s\n", name, Network)))
				} else {
					discrete, _ := service.GetFieldBool("discrete")
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
		if !resolver.Status(&model.Params{Domain: c.Domain}) {
			resolver.Configure(&model.Params{Domain: c.Domain})
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
		name, _ := service.GetFieldString("name")
		url, _ := service.GetFieldString("url")
		if s, _ := service.Status(); s && url != "" {
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
