package library

import (
	"fmt"
	"strings"

	"github.com/pygmystack/pygmy/service/interface/docker"
)

// Update will update the the images for all configured services.
func Update(c Config) {

	// Import the configuration.
	Setup(&c)

	// Loop over services.
	for s := range c.Services {

		// Pull the image.
		service := c.Services[s]
		purpose, _ := service.GetFieldString("purpose")
		var result string
		var err error
		if purpose == "" || purpose == "sshagent" {
			result, err = docker.DockerPull(service.Config.Image)
			if err == nil {
				fmt.Println(result)
			}
		}

		// If the service is running, restart it.
		if s, _ := service.Status(); s && !strings.Contains(result, "is up to date") {
			var e error
			e = service.Stop()
			if e != nil {
				fmt.Println(e)
			}
			if s, _ := service.Status(); !s {
				e = service.Start()
				if e != nil {
					fmt.Println(e)
				}
			}
		}
	}

}
