// Copyright © 2019 Karl Hepworth <Karl.Hepworth@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/fubarhouse/pygmy/v1/service/library"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Example: "pygmy up",
	Short: "Bring up pygmy services (dnsmasq, haproxy, mailhog, resolv, ssh-agent)",
	Long: `Launch Pygmy - a set of containers and a resolver with very specific
configurations designed for use with Amazee.io local development.

It includes dnsmasq, haproxy, mailhog, resolv and ssh-agent.`,
	Run: func(cmd *cobra.Command, args []string) {

		c.Key, _ = cmd.Flags().GetString("key")
		c.SkipKey, _ = cmd.Flags().GetBool("no-addkey")
		c.SkipResolver, _ = cmd.Flags().GetBool("no-resolver")

		library.Up(c)

	},
}

func init() {

	homedir, _ := homedir.Dir()
	keypath := fmt.Sprintf("%v%v.ssh%vid_rsa", homedir, string(os.PathSeparator), string(os.PathSeparator))

	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringP("key", "", keypath, "Path of SSH key to add")
	upCmd.Flags().BoolP("no-addkey", "", false, "Skip adding the SSH key")
	upCmd.Flags().BoolP("no-resolver", "", false, "Skip adding or removing the Resolver")

}
