package library

import (
	"fmt"
	"os"
	"runtime"

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

		for _, Container := range c.Services {
			if Container.Group == "addkeys" {
				if runtime.GOOS == "windows" {
					Container.Config.Cmd = []string{"ssh-add", "/key"}
					Container.HostConfig.Binds = append(Container.HostConfig.Binds, fmt.Sprintf("%v:/key", key))
				} else {
					Container.Config.Cmd = []string{"ssh-add", key}
					Container.HostConfig.Binds = append(Container.HostConfig.Binds, fmt.Sprintf("%v:%v", key, key))
				}
				// TODO: Need to figure out why this specifically isn't working. This should resolve #33.
				Container.Start()
			}
		}

	} else {
		fmt.Printf("Already added key file %v.\n", key)
	}
}