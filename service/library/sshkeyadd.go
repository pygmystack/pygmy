package library

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	. "github.com/logrusorgru/aurora"

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
				color.Print(Green(fmt.Sprintf("Validation success for SSH key %v\n", key)))
			} else {
				if err.Error() == "ssh: this private key is passphrase protected" {
					color.Print(Green(fmt.Sprintf("Validation success for protected SSH key %v\n", key)))
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

			if err := Container.Create(); err != nil {
				_ = Container.Remove()
				return err
			}
			if err := Container.Start(); err != nil {
				_ = Container.Remove()
				return err
			}

			interactive, _ := Container.GetFieldBool("interactive")
			if !interactive {
				l, _ := Container.DockerLogs()
				handled := false
				// We need tighter control on the output of this container...
				for _, line := range strings.Split(string(l), "\n") {
					if strings.Contains(line, "Identity added:") {
						handled = true
						color.Print(Green(fmt.Sprintf("Successfully added SSH key %v to agent\n", key)))
					}
					if strings.Contains(line, "Enter passphrase for") {
						handled = true
						color.Print(Yellow("Warning: Passphrase protected SSH keys can only be added in interactive mode, the key will not be added.\n"))
					}
				}

				// Logs didn't contain known messages, log all in case of error.
				if !handled {
					color.Print(Red("Unknown error while adding SSH key:\n"))
					for _, line := range strings.Split(string(l), "\n") {
						fmt.Println(line)
					}
				}
			}

			_ = Container.Remove()

		}

	}
	return nil
}
