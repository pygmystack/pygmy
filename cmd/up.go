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

	"github.com/mitchellh/go-homedir"
	"github.com/pygmystack/pygmy/external/docker/commands"
	"github.com/pygmystack/pygmy/external/docker/setup"
	"github.com/pygmystack/pygmy/internal/utils/cert"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"start"},
	Example: "pygmy up",
	Short:   "Bring up pygmy services (dnsmasq, haproxy, mailhog, resolv, ssh-agent)",
	Long: `Launch Pygmy - a set of containers and a resolver with very specific
configurations designed for use with Amazee.io local development.
It includes dnsmasq, haproxy, mailhog, resolv and ssh-agent.`,
	Run: func(cmd *cobra.Command, args []string) {
		Key, _ := cmd.Flags().GetString("key")
		NoKey, _ := cmd.Flags().GetBool("no-addkey")
		noResolv, _ := cmd.Flags().GetBool("no-resolver")
		c.TLSCertPath, _ = cmd.Flags().GetString("tls-cert")

		if noResolv {
			c.ResolversDisabled = true
		}
		if NoKey {
			c.Keys = []setup.Key{}
		} else {
			keyExistsInConfig := false
			for _, key := range c.Keys {
				if key.Path == Key {
					keyExistsInConfig = true
				}
			}
			if !keyExistsInConfig {
				thisKey := setup.Key{
					Path: Key,
				}
				c.Keys = append(c.Keys, thisKey)
			}
		}

		err := commands.Up(c)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {

	homedir, _ := homedir.Dir()
	keypath := fmt.Sprintf("%v%v.ssh%vid_rsa", homedir, string(os.PathSeparator), string(os.PathSeparator))
	tlsCertdefault := cert.GetDefaultCertPath()

	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringP("key", "", keypath, "Path of SSH key to add")
	upCmd.Flags().BoolP("no-addkey", "", false, "Skip adding the SSH key")
	upCmd.Flags().BoolP("no-resolver", "", false, "Skip adding or removing the Resolver")
	upCmd.Flags().StringP("tls-cert", "", tlsCertdefault, "Path to TLS certificate to use with the Pygmy haproxy")
}
