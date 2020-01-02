package library

import "github.com/fubarhouse/pygmy-go/service/resolv"

func Down(c Config) {

	Setup(&c)

	for _, Service := range c.Services {
		if !Service.Disabled {
			Service.Stop()
		}
	}

	for _, resolver := range c.Resolvers {
		resolv.New(resolver).Clean()
	}
}