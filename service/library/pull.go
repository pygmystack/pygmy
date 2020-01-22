package library

import (
	"github.com/fubarhouse/pygmy-go/service/interface"
)

func Pull(c Config) {

	Setup(&c)

	for _, Service := range c.Services {

		model.DockerPull(Service.Config.Image)

	}
}
