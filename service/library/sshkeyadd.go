package library

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/fubarhouse/pygmy-go/service/ssh/agent"
)

func SshKeyAdd(c Config, key string) ([]byte, error) {

	if c.SkipKey {
		return []byte{}, nil
	}

	Setup(&c)

	if key != "" {
		if _, err := os.Stat(key); err != nil {
			fmt.Printf("%v\n", err)
			return []byte{}, err
		}
	}

	var b []byte
	var e error

	if !agent.Search(key) {

		for _, Container := range c.Services {
			purpose, _ := Container.GetFieldString("purpose")
			if purpose == "addkeys" {
				if runtime.GOOS == "windows" {
					Container.Config.Cmd = []string{"ssh-add", "/key"}
					Container.HostConfig.Binds = append(Container.HostConfig.Binds, fmt.Sprintf("%v:/key", key))
				} else {
					Container.Config.Cmd = []string{"ssh-add", key}
					Container.HostConfig.Binds = append(Container.HostConfig.Binds, fmt.Sprintf("%v:%v", key, key))
				}

				b, e = Container.Start()
			}
		}

	} else {
		e = errors.New(fmt.Sprintf("Already added key file %v.\n", key))
		return b, e
	}
	return b, e
}
