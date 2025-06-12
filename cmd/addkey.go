// Copyright © 2019 Karl Hepworth <Karl.Hepworth@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package cmd

import (
	"fmt"
	"github.com/pygmystack/pygmy/external/docker/setup"
	"strings"

	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/pygmystack/pygmy/external/docker/commands"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/containers"
	"github.com/pygmystack/pygmy/internal/utils/color"
)

// addkeyCmd is the SSH key add command.
var addkeyCmd = &cobra.Command{
	Use:     "addkey",
	Example: "pygmy addkey --key ~/.ssh/id_rsa",
	Short:   "Add/re-add an SSH key to the agent",
	Long:    `Add or re-add an SSH key to Pygmy's SSH Agent by specifying the path to the private key.`,
	Run: func(cmd *cobra.Command, args []string) {

		cli, ctx, err := internals.NewClient()
		if err != nil {
			fmt.Println(err)
		}

		Key, _ := cmd.Flags().GetString("key")
		var Keys []setup.Key

		if Key != "" {
			thisKey := setup.Key{
				Path: Key,
			}
			Keys = append(Keys, thisKey)
		} else {
			if len(Keys) == 0 {
				setup.Setup(ctx, cli, &c)
				Keys = c.Keys
			}
		}

		for _, k := range Keys {
			if e := commands.SshKeyAdd(c, k.Path); e != nil {
				color.Print(Red(fmt.Sprintf("%v\n", e)))
			}
		}

		for _, s := range c.SortedServices {
			service := c.Services[s]
			purpose, _ := service.GetFieldString(ctx, cli, "purpose")
			if purpose == "sshagent" {
				name, _ := service.GetFieldString(ctx, cli, "name")
				d, _ := containers.Exec(ctx, cli, name, "ssh-add -l")
				if strings.Contains(string(d), "The agent has no identities.") {
					fmt.Println(Red("Agent has no identities, the key could not be added."))
					fmt.Println(Red("Start the SSH Agent, add the SSH key and try again."))
				} else {
					fmt.Println(string(d))
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(addkeyCmd)
	addkeyCmd.Flags().StringP("key", "k", "", "Path of SSH key to add")
}
