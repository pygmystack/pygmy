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
	"github.com/fubarhouse/pygmy/service/library"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Example: "pygmy down",
	Short: "Stop and remove all pygmy services",
	Long: `Check if any pygmy containers are running and removes
then if they are, it will not attempt to remove any
services which are not running.`,
	Run: func(cmd *cobra.Command, args []string) {

		library.Down(c)

	},
}

func init() {

	rootCmd.AddCommand(downCmd)

}
