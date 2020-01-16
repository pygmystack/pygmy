package library

import (
	"github.com/fubarhouse/pygmy-go/service/resolv"
)

func Down(c Config) {

	Setup(&c)

	for _, Service := range c.Services {
		disabled, _ := Service.GetFieldBool("disabled")
		if !disabled {
			Service.Stop()
		}
	}

	for _, resolver := range c.Resolvers {
		resolv.New(resolver).Clean()
	}
}
