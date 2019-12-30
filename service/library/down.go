package library

import "github.com/fubarhouse/pygmy-go/service/resolv"

// Down provides business logic for the `down` command.
func Down(c Config) {

	Setup(&c)

	for _, Service := range c.Services {
		if !Service.Disabled {
			Service.Stop()
		}
	}

	if !c.SkipResolver {
		for _, resolver := range c.Resolvers {
			resolv.New(resolver).Clean()
		}
	}
}