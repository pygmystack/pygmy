package library

import (
	"fmt"

	"github.com/fubarhouse/pygmy-go/service/interface/docker"
	. "github.com/logrusorgru/aurora"
)

// Down will bring pygmy down safely
func Down(c Config) {

	Setup(&c)
	NetworksToClean := []string{}

	for _, Service := range c.Services {
		enabled, _ := Service.GetFieldBool("enable")
		if enabled {
			e := Service.Stop()
			if e != nil {
				fmt.Println(e)
			}
			if s, _ := Service.GetFieldString("network"); s != "" {
				NetworksToClean = append(NetworksToClean, s)
			}
		}
	}

	for _, network := range c.Networks {
		NetworksToClean = append(NetworksToClean, network.Name)
	}

	for _, network := range unique(NetworksToClean) {
		e := docker.DockerNetworkRemove(network)
		if e != nil {
			fmt.Println(e)
		}
		if s, _ := docker.DockerNetworkStatus(network); !s {
			fmt.Printf(Sprintf(Green("Successfully removed network %v\n"), Green(network)))
		} else {
			fmt.Printf(Sprintf(Red("Network %v was not removed"), Red(network)))
		}
	}

	for _, resolver := range c.Resolvers {
		resolver.Clean()
	}
}
