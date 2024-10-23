package commands

import (
	"fmt"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
	"os"
	"runtime"
	"strings"

	. "github.com/logrusorgru/aurora"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/containers"
	"github.com/pygmystack/pygmy/internal/service/docker/ssh/agent"
	"github.com/pygmystack/pygmy/internal/utils/color"
)

// SshKeyAdd will add a given key to the ssh agent.
func SshKeyAdd(c Config, key string) error {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return err
	}

	Setup(ctx, cli, &c)

	if key != "" {
		if _, err := os.Stat(key); err != nil {
			fmt.Printf("%v\n", err)
			return err
		}
	} else {
		return nil
	}

	for _, Container := range c.Services {
		purpose, _ := Container.GetFieldString(ctx, cli, "purpose")
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
				Container.Config.Cmd = []string{"windows-key-add", "/key"}
				Container.HostConfig.Binds = append(Container.HostConfig.Binds, fmt.Sprintf("%v:/key", key))
			} else {
				Container.Config.Cmd = []string{"ssh-add", key}
				Container.HostConfig.Binds = append(Container.HostConfig.Binds, fmt.Sprintf("%v:%v", key, key))
			}

			if err := Container.Create(ctx, cli); err != nil {
				_ = Container.Remove(ctx, cli)
				return err
			}
			if err := Container.Start(ctx, cli); err != nil {
				_ = Container.Remove(ctx, cli)
				return err
			}

			interactive, _ := Container.GetFieldBool(ctx, cli, "interactive")
			name, _ := Container.GetFieldString(ctx, cli, "name")
			if !interactive {
				l, _ := containers.Logs(ctx, cli, name)
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

			_ = Container.Remove(ctx, cli)

		}

	}
	return nil
}
