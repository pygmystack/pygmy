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
	"github.com/pygmystack/pygmy/external/commands"
	"github.com/spf13/cobra"
)

var jsonOutput bool

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:     "status",
	Example: "pygmy status",
	Short:   "Report status of the pygmy services",
	Long: `Loop through all of pygmy's services and identify the present state.
This includes the docker services, the resolver and SSH key status`,
	Run: func(cmd *cobra.Command, args []string) {

		if jsonOutput {
			c.JSONFormat = true
		}
		commands.Status(c)

	},
}

func init() {

	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVarP(&jsonOutput, "json", "", false, "Output status in JSON format")

}
