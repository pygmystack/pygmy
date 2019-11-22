package library

import (
	"fmt"
	"os"

	"github.com/fubarhouse/pygmy-go/service/ssh/agent"
)

func SshKeyAdd(c Config, key string) {

	if c.SkipKey {
		return
	}

	Setup(&c)

	if key != "" {
		if _, err := os.Stat(key); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}

	if !agent.Search(key) {

		if key != "" {
			c.Key = key
		}

		for _, Container := range c.Services {
			if Container.Group == "addkey" {
				Container.Start()
			}
		}

	} else {
		fmt.Printf("Already added key file %v.\n", c.Key)
	}
}