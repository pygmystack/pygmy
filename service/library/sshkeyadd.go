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
			vOne, _ := Container.TagGet("addkeys")
			if vOne == "pygmy.addkeys" {
				if runtime.GOOS == "windows" {
					Container.Config.Cmd = []string{"ssh-add", "/key"}
					Container.HostConfig.Binds = append(Container.HostConfig.Binds, fmt.Sprintf("%v:/key", key))
				} else {
					Container.Config.Cmd = []string{"ssh-add", key}
					Container.HostConfig.Binds = append(Container.HostConfig.Binds, fmt.Sprintf("%v:%v", key, key))
				}
				Container.Start()
			}
		}

		for _, Container := range c.Services {
			vOne, _ := Container.TagGet("addkeys")
			if vOne == "pygmy.addkeys" {
				Container.Start()
			}
		}

	} else {
		fmt.Printf("Already added key file %v.\n", key)
	}
}
