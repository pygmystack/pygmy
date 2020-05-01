package library

import (
	"fmt"
	"github.com/fubarhouse/pygmy-go/service/interface"
)

func Pull(c Config) {

	Setup(&c)

	for _, Service := range c.Services {

		_, e := model.DockerPull(Service.Config.Image)
		if e != nil {
			fmt.Print(e)
		}

	}
}
