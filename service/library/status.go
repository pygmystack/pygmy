package library

import (
	"fmt"
	"strings"

	"github.com/fubarhouse/pygmy-go/service/haproxy_connector"
	model "github.com/fubarhouse/pygmy-go/service/interface"
	"github.com/fubarhouse/pygmy-go/service/resolv"
	"github.com/fubarhouse/pygmy-go/service/test_url"
)

func Status(c Config) {

	Setup(&c)
	checks := DryRun(&c)

	if len(checks) > 0 {
		for _, check := range checks {
			fmt.Println(check.Message)
		}
		fmt.Println()
	}

	Containers, _ := model.DockerContainerList()
	for _, Container := range Containers {
		if Container.Labels["pygmy"] == "pygmy" {
			Service := c.Services[strings.Trim(Container.Names[0], "/")]
			if s, _ := Service.Status(); s {
				name, _ := Service.GetFieldString("name")
				disabled, _ := Service.GetFieldBool("disabled")
				discrete, _ := Service.GetFieldBool("discrete")
				if name != "" {
					if !disabled && !discrete && name != "" {
						if s, _ := Service.Status(); s {
							fmt.Printf("[*] %v: Running as container %v\n", name, name)
						} else {
							fmt.Printf("[ ] %v is not running\n", name)
						}
					}
				}
			}
		}
	}

	for _, Service := range c.Services {
		if s, _ := Service.Status(); !s {
			name, _ := Service.GetFieldString("name")
			fmt.Printf("[ ] %v is not running\n", name)
		}
	}

	for _, Network := range c.Networks {
		netStat, _ := NetworkStatus(Network.Name)
		if netStat {
			for _, Container := range Network.Containers {
				if s, _ := haproxy_connector.Connected(Container, Network.Name); s {
					fmt.Printf("[*] %v is connected to network %v\n", Container, Network.Name)
				} else {
					fmt.Printf("[ ] %v is not connected to network %v\n", Container, Network.Name)
				}
			}
		}
		if _, e := NetworkStatus(Network.Name); e == nil {
			fmt.Printf("[*] %v network has been created\n", Network.Name)
		} else {
			fmt.Printf("[ ] %v network has not been created\n", Network.Name)
		}
	}

	for _, resolver := range c.Resolvers {
		r := resolv.New(resolv.Resolv{Name: resolver.Name, Data: resolver.Data, Folder: resolver.Folder, File: resolver.File})
		if s := r.Status(); s {
			fmt.Printf("[*] Resolv %v is properly connected\n", resolver.Name)
		} else {
			fmt.Printf("[ ] Resolv %v is not properly connected\n", resolver.Name)
		}
	}

	for _, volume := range c.Volumes {
		if s, _ := model.DockerVolumeExists(volume); s {
			fmt.Printf("[*] Volume %v has been created\n", volume.Name)
		} else {
			fmt.Printf("[ ] Volume %v has not been created\n", volume.Name)
		}
	}

	for _, Container := range c.Services {
		purpose, _ := Container.GetFieldString("purpose")
		//output, _ := Container.GetFieldBool("output")
		if purpose == "showkeys" {
			// TODO re-fix
			//Container.Output = true
			Container.Start()
		}
	}

	for _, Container := range c.Services {
		Status, _ := Container.Status()
		name, _ := Container.GetFieldString("name")
		url, _ := Container.GetFieldString("url")
		if url != "" && Status {
			test_url.Validate(url)
			if r := test_url.Validate(url); r {
				fmt.Printf(" - %v (%v)\n", url, name)
			} else {
				fmt.Printf(" ! %v (%v)\n", url, name)
			}
		}
	}

}
