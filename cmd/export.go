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

	"github.com/pygmystack/pygmy/external/docker/commands"
)

var exportPath string

// exportCmd represents the status command
var exportCmd = &cobra.Command{
	Use:     "export",
	Example: "pygmy export --config /path/to/input --output /path/to/output",
	Short:   "Export validated configuration to a given path",
	Long:    `Export configuration which has validated into a specified path`,
	Run: func(cmd *cobra.Command, args []string) {

		err := commands.Export(c, exportPath)
		if err != nil {
			fmt.Println(err)
		}

	},
}

func init() {

	rootCmd.AddCommand(exportCmd)

	homedir, _ := homedir.Dir()
	exportPath = fmt.Sprintf("%v%v.pygmy.yml", homedir, string(os.PathSeparator))

	exportCmd.Flags().StringVarP(&exportPath, "output", "o", exportPath, "Path to exported configuration to be written to")

}
