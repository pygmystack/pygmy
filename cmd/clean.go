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

// stopCmd represents the stop command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Example: "pygmy clean",
	Short: "Stop and remove all pygmy services regardless of state",
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
