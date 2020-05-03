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
	"github.com/fubarhouse/pygmy-go/service/library"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var cleanCmd = &cobra.Command{
	Use:     "clean",
	Example: "pygmy-go clean",
	Short:   "Stop and remove all pygmy services regardless of state",
	Long: `Useful for debugging or system cleaning, this command will
remove all pygmy containers but leave the images in-tact.

This command does not check if the containers are running
because other checks do for speed convenience.`,
	Run: func(cmd *cobra.Command, args []string) {

		library.Clean(c)

	},
}

func init() {

	rootCmd.AddCommand(cleanCmd)

}
