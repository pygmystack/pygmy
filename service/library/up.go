package library

import (
	"os"
	"fmt"

	"github.com/fubarhouse/pygmy-go/service/haproxy_connector"
	model "github.com/fubarhouse/pygmy-go/service/interface"
	"github.com/fubarhouse/pygmy-go/service/network"
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
			if err != nil {
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
	for _, service := range c.SortedServices {
		s := c.Services[service]
		if !s.Disabled && s.Group != "addkeys" && s.Group != "showkeys" {
			s.Start()
		}
	}

	for Network, Containers := range c.Networks {
		netStat, _ := network.Status(Network)
		if !netStat {
			network.Create(Network)
		}
		for _, Container := range Containers {
			if s, _ := haproxy_connector.Connected(Container, Network); !s {
				haproxy_connector.Connect(Container, Network)
				if s, _ := haproxy_connector.Connected(Container, Network); s {
					fmt.Printf("Successfully connected %v to %v\n", Container, Network)
				} else {
					fmt.Printf("Could not connect %v to %v\n", Container, Network)
				}
			} else {
				fmt.Printf("Already connected %v to %v\n", Container, Network)
			}
		}
	}

	if !c.SkipResolver {
		for _, resolver := range c.Resolvers {
			resolv.New(resolver).Configure()
		}
	}

	if !c.SkipKey {

		for _, key := range c.Keys {
			SshKeyAdd(c, key)
		}

	}

	for _, service := range c.Services {
		if s, _ := service.Status(); s && service.URL != "" {
			test_url.Validate(service.URL)
			if r := test_url.Validate(service.URL); r {
				fmt.Printf(" - %v (%v)\n", service.URL, service.Name)
			} else {
				fmt.Printf(" ! %v (%v)\n", service.URL, service.Name)
			}
		}
	}
}
