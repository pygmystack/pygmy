package library

import (
	"fmt"
	"os"
	"strings"

	"github.com/fubarhouse/pygmy-go/service/endpoint"
	model "github.com/fubarhouse/pygmy-go/service/interface"
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
		if s, _ := model.DockerVolumeExists(volume); !s {
			_, err := model.DockerVolumeCreate(volume)
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
		output, _ := service.GetFieldBool("output")

		// Do not show or add keys:
		if !disabled && purpose != "addkeys" && purpose != "showkeys" {

			// Here we will immitate the docker command by
			// pulling the image if it's not in the daemon.
			images, _ := model.DockerImageList()
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
			// When running 'docker run', it will pull the image.
			// For UX it makes sense we do this here.
			if !imageFound {
				if err := model.DockerPull(service.Config.Image); err != nil {
					continue
				}
			}

			o, _ := service.Start()
			if output && string(o) != "" {
				fmt.Println(string(o))
			}
		}

		// If one or more agent was found:
		if purpose == "sshagent" {
			agentPresent = true
		}
	}

	for _, Network := range c.Networks {
		netStat, _ := NetworkStatus(Network.Name)
		if !netStat {
			if err := NetworkCreate(Network); err == nil {
				fmt.Printf("Successfully created network %v\n", Network.Name)
			} else {
				fmt.Printf("Could not create network %v\n", Network.Name)
			}

		}
		for _, Container := range Network.Containers {
			if s, _ := model.DockerNetworkConnected(Network, Container.Name); !s {
				if s := model.DockerNetworkConnect(Network, Container.Name); s == nil {
					fmt.Printf("Successfully connected %v to %v\n", Container.Name, Network.Name)
				} else {
					fmt.Printf("Could not connect %v to %v\n", Container.Name, Network.Name)
				}
			} else {
				fmt.Printf("Already connected %v to %v\n", Container.Name, Network.Name)
			}
		}
	}

	for _, resolver := range c.Resolvers {
		if resolver.Enabled {
			resolver.Configure()
		}
	}

	// Add ssh-keys to the agent
	if agentPresent {
		i := 1
		for _, v := range c.Keys {
			out, err := SshKeyAdd(c, v, i)
			if err != nil {
				fmt.Println(err)
			} else if string(out) != "" {
				fmt.Println(strings.Trim(string(out), "\n"))
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
