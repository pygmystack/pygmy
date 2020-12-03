package library

import (
	"fmt"
	"os"
	"strings"

	"github.com/fubarhouse/pygmy-go/service/endpoint"
	"github.com/fubarhouse/pygmy-go/service/interface/cri/docker"
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
		if s, _ := docker.VolumeExists(volume); !s {
			_, err := docker.VolumeCreate(volume)
			if err == nil {
				fmt.Printf("Created volume %v\n", volume.Name)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("Already created volume %v\n", volume.Name)
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

		// Do not show or add keys:
		if enabled && purpose != "addkeys" && purpose != "showkeys" {

			// Here we will immitate the command by
			// pulling the image if it's not in the daemon.
			images, _ := docker.ImageList()
			imageFound := false
			for _, image := range images {
				for _, digest := range image.RepoDigests {
					d := strings.Trim(strings.SplitAfter(digest, "@")[0], "@")
					if strings.Contains(service.Config.Image, d) {
						imageFound = true
					}
				}
			}

			// The image wasn't found.
			// When running the container, it will pull the image.
			// For UX it makes sense we do this here.
			if !imageFound {
				if _, err := docker.ImagePull(service.Config.Image); err != nil {
					continue
				}
			}

			e := service.Start()
			if e != nil {
				fmt.Println(e)
			}
		}

		// If one or more agent was found:
		if purpose == "sshagent" {
			agentPresent = true
		}
	}

	// Network(s) creation
	for _, Network := range c.Networks {
		if Network.Name != "" {
			netVal, _ := docker.NetworkStatus(Network.Name)
			if !netVal {
				if err := NetworkCreate(Network); err == nil {
					fmt.Printf("Successfully created network %v\n", Network.Name)
				} else {
					fmt.Printf("Could not create network %v\n", Network.Name)
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
			if s, _ := docker.NetworkConnected(Network, name); !s {
				if s := NetworkConnect(Network, name); s == nil {
					fmt.Printf("Successfully connected %v to %v\n", name, Network)
				} else {
					discrete, _ := service.GetFieldBool("discrete")
					if !discrete {
						fmt.Printf("Could not connect %v to %v\n", name, Network)
					}
				}
			} else {
				fmt.Printf("Already connected %v to %v\n", name, Network)
			}
		}
	}

	for _, resolver := range c.Resolvers {
		if !resolver.Status() {
			resolver.Configure()
		}
	}

	// Add ssh-keys to the agent
	if agentPresent {
		i := 1
		for _, v := range c.Keys {
			err := SshKeyAdd(c, v, i)
			if err != nil {
				fmt.Println(err)
			}
			i++
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
}
