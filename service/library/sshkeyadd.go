package library

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fubarhouse/pygmy-go/service/ssh/agent"
)

// SshKeyAdd will add a given key to the ssh agent.
func SshKeyAdd(c Config, key string, index int) error {

	Setup(&c)

	if key != "" {
		if _, err := os.Stat(key); err != nil {
			fmt.Printf("%v\n", err)
			return err
		}
	}

	var e error

	for _, Container := range c.Services {
		purpose, _ := Container.GetFieldString("purpose")
		if purpose == "addkeys" {
			if !agent.Search(Container, key) {
				if runtime.GOOS == "windows" {
					Container.Config.Cmd = []string{"ssh-add", "/key"}
					Container.HostConfig.Binds = append(Container.HostConfig.Binds, fmt.Sprintf("%v:/key", key))
				} else {
					Container.Config.Cmd = []string{"ssh-add", key}
					Container.HostConfig.Binds = append(Container.HostConfig.Binds, fmt.Sprintf("%v:%v", key, key))
				}

				if e := Container.Create(); e != nil {
					return e
				}
				// THIS.
				if e := Container.Start(); e != nil {
					return e
				}
				l, _ := Container.DockerLogs()
				Container.Remove()

				// We need tighter control on the output of this container...
				for _, line := range strings.Split(string(l), "\n") {
					if strings.Contains(line, "Identity added:") {
						fmt.Println(line)
					}
				}

			}

		}
	}
	return e
}
