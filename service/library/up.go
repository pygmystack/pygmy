package library

import (
	"fmt"
	"os"

	"github.com/fubarhouse/pygmy-go/service/haproxy_connector"
	model "github.com/fubarhouse/pygmy-go/service/interface"
	"github.com/fubarhouse/pygmy-go/service/resolv"
	"github.com/fubarhouse/pygmy-go/service/test_url"
)

func Up(c Config) {

	Setup(&c)
	checks := DryRun(&c)

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
			fmt.Printf("Already created volume %v\n", volume)
		}
	}

	// Maps are... bad for predictable sequencing.
	// Look over the sorted slice and start them in
	// alphabetical order - so that one can configure
	// an ssh-agent like amazeeio-ssh-agent.
	for _, s := range c.SortedServices {
		service := c.Services[s]
		disabled, _ := service.GetFieldBool("disabled")
		purpose, _ := service.GetFieldString("purpose")
		if !disabled && purpose != "addkeys" && purpose != "showkeys" {
			service.Start()
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
			if s, _ := haproxy_connector.Connected(Container.Name, Network.Name); !s {
				haproxy_connector.Connect(Container.Name, Network.Name)
				if s, _ := haproxy_connector.Connected(Container.Name, Network.Name); s {
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
		this := resolv.New(resolver)
		if !this.Disabled {
			this.Configure()
		}
	}

	if !c.SkipKey {

		for _, key := range c.Keys {
			SshKeyAdd(c, key)
		}

	}

	for _, service := range c.Services {
		name, _ := service.GetFieldString("name")
		url, _ := service.GetFieldString("url")
		if s, _ := service.Status(); s && url != "" {
			test_url.Validate(url)
			if r := test_url.Validate(url); r {
				fmt.Printf(" - %v (%v)\n", url, name)
			} else {
				fmt.Printf(" ! %v (%v)\n", url, name)
			}
		}
	}
}
