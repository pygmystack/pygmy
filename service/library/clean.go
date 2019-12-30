package library

import (
	"fmt"

	model "github.com/fubarhouse/pygmy-go/service/model"
	"github.com/fubarhouse/pygmy-go/service/resolv"
)

// Clean provides business logic for the `clean` command.
func Clean(c Config) {

	Setup(&c)
	Containers, _ := model.DockerContainerList()

	for _, Container := range Containers {
		if l := Container.Labels["pygmy"]; l == "pygmy" {

			err := model.DockerKill(Container.ID)
			if err == nil {
				fmt.Printf("Successfully killed  %v.\n", Container.Names[0])
			}

			err = model.DockerRemove(Container.ID)
			if err == nil {
				fmt.Printf("Successfully removed %v.\n", Container.Names[0])
			}
		}
	}

	for _, resolver := range c.Resolvers {
		resolv.New(resolver).Clean()
	}
}
