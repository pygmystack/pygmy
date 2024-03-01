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
	"github.com/spf13/cobra"

	"github.com/pygmystack/pygmy/service/library"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:     "restart",
	Example: "pygmy restart",
	Short:   "Restart all pygmy containers.",
	Long:    `This command will trigger the Down and Up commands`,
	Run: func(cmd *cobra.Command, args []string) {

		Key, _ := cmd.Flags().GetString("key")
		NoKey, _ := cmd.Flags().GetBool("no-addkey")

		if NoKey {
			c.Keys = []library.Key{}
		} else {
			FoundKey := false
			for _, v := range c.Keys {
				if v.Path == Key {
					FoundKey = true
				}
			}

			if !FoundKey {
				thisKey := library.Key{
					Path: Key,
				}
				c.Keys = append(c.Keys, thisKey)
			}
		}

		library.Restart(c)

	},
}

func init() {

	homedir, _ := homedir.Dir()
	keypath := fmt.Sprintf("%v%v.ssh%vid_rsa", homedir, string(os.PathSeparator), string(os.PathSeparator))

	rootCmd.AddCommand(restartCmd)
	restartCmd.Flags().StringP("key", "", keypath, "Path of SSH key to add")
	restartCmd.Flags().BoolP("no-addkey", "", false, "Skip adding the SSH key")
	restartCmd.Flags().BoolP("no-resolver", "", false, "Skip adding or removing the Resolver")
}
