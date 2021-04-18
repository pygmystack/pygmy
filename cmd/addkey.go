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
	"os"

	"github.com/fubarhouse/pygmy-go/service/interface/docker"
	"github.com/fubarhouse/pygmy-go/service/library"
	"github.com/spf13/cobra"
)

// addkeyCmd is the SSH key add command.
var addkeyCmd = &cobra.Command{
	Use:     "addkey",
	Example: "pygmy addkey --key ~/.ssh/id_rsa",
	Short:   "Add/re-add an SSH key to the agent",
	Long:    `Add or re-add an SSH key to Pygmy's SSH Agent by specifying the path to the private key.`,
	Run: func(cmd *cobra.Command, args []string) {

		Key, _ := cmd.Flags().GetString("key")
		Keys := []string{}

		if Key != "" {
			Keys = append(Keys, Key)
		} else {
			if _, err := os.Stat(os.Args[len(os.Args) - 1]); err == os.ErrExist {
				Keys = append(c.Keys, os.Args[len(os.Args) - 1])
			}
			if len(Keys) == 0 {
				library.Setup(&c)
				Keys = c.Keys
			}
		}

		for _, k := range Keys {
			if e := library.SshKeyAdd(c, k); e != nil {
				fmt.Println(e)
			}
		}

		for _, s := range c.SortedServices {
			service := c.Services[s]
			purpose, _ := service.GetFieldString("purpose")
			if purpose == "sshagent" {
				name, _ := service.GetFieldString("name")
				d, _ := docker.DockerExec(name, "ssh-add -l")
				fmt.Println(string(d))
			}
		}

	},
}

func init() {

	rootCmd.AddCommand(addkeyCmd)
	addkeyCmd.Flags().StringP("key", "", "", "Path of SSH key to add")

}
