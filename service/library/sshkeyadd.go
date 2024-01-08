package library

import (
	"fmt"
	"os"
	"runtime"

	"github.com/logrusorgru/aurora"

	"github.com/pygmystack/pygmy/service/color"
	"github.com/pygmystack/pygmy/service/ssh/agent"
)

// SshKeyAdd will add a given key to the ssh agent.
func SshKeyAdd(c Config, key string) error {

	Setup(&c)

	if key != "" {
		if _, err := os.Stat(key); err != nil {
			fmt.Printf("%v\n", err)
			return err
		}
	} else {
		return nil
	}

	for _, Container := range c.Services {
		purpose, _ := Container.GetFieldString("purpose")
		if purpose == "addkeys" {

			// Validate SSH Key before adding.
			valid, err := agent.Validate(key)
			if valid {
				color.Print(aurora.Green(fmt.Sprintf("Validation success for SSH key %v\n", key)))
			} else {
				if err.Error() == "ssh: this private key is passphrase protected" {
					color.Print(aurora.Green(fmt.Sprintf("Validation success for protected SSH key %v\n", key)))
				}
				if err.Error() == "ssh: no key found" {
					return fmt.Errorf(fmt.Sprintf("[ ] Validation failure for SSH key %v\n", key))
				}
			}

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
			if e := Container.Start(); e != nil {
				return e
			}
			l, _ := Container.DockerLogs()
			_ = Container.Remove()

			// We need tighter control on the output of this container...
			for _, line := range strings.Split(string(l), "\n") {
				if strings.Contains(line, "Identity added:") {
					color.Print(aurora.Green(fmt.Sprintf("Successfully added SSH key %v to agent\n", key)))
				}
				if strings.Contains(line, "Enter passphrase for") {
					color.Print(aurora.Yellow("Warning: Passphrase protected SSH keys are not currently supported, the key will not be added.\n"))
				}
			}

		}

	}
	return nil
}
