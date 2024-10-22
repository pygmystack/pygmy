package commands

import (
	"fmt"
	"strings"

	runtimeimages "github.com/pygmystack/pygmy/internal/runtime/docker/docker/images"
)

// Update will update the images for all configured services.
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
			result, err = runtimeimages.Pull(service.Config.Image)
			if err == nil {
				fmt.Println(result)
			} else {
				fmt.Println(err)
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

	images, _ := runtimeimages.List()
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if strings.Contains(tag, "uselagoon") {
				result, err := runtimeimages.Pull(tag)
				if err == nil {
					fmt.Println(result)
				}
			}
		}
	}
}
