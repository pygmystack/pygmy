package commands

import (
	"fmt"
	"strings"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
	runtimeimages "github.com/pygmystack/pygmy/internal/runtime/docker/internals/images"
)

// Update will update the images for all configured services.
func Update(c Config) error {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return err
	}

	// Import the configuration.
	Setup(ctx, cli, &c)

	// Loop over services.
	for s := range c.Services {

		// Pull the image.
		service := c.Services[s]
		purpose, _ := service.GetFieldString(ctx, cli, "purpose")
		var result string
		var err error
		if purpose == "" || purpose == "sshagent" {
			result, err = runtimeimages.Pull(ctx, cli, service.Config.Image)
			if err == nil {
				fmt.Println(result)
			} else {
				fmt.Println(err)
			}
		}

		// If the service is running, restart it.
		if s, _ := service.Status(ctx, cli); s && !strings.Contains(result, "is up to date") {
			var e error
			e = service.Stop(ctx, cli)
			if e != nil {
				fmt.Println(e)
			}
			if s, _ := service.Status(ctx, cli); !s {
				e = service.Start(ctx, cli)
				if e != nil {
					fmt.Println(e)
				}
			}
		}
	}

	images, _ := runtimeimages.List(ctx, cli)
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if strings.Contains(tag, "uselagoon") {
				result, err := runtimeimages.Pull(ctx, cli, tag)
				if err == nil {
					fmt.Println(result)
				}
			}
		}
	}

	return nil
}
