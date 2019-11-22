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

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Example: "pygmy update",
	Short: "Pulls Docker Images and recreates the Containers",
	Long: `Pull all images Pygmy uses, as well as any images containing
the string 'amazeeio', which encompasses all lagoon images.`,
	Run: func(cmd *cobra.Command, args []string) {

		library.Update(c)

	},
}

func init() {

	rootCmd.AddCommand(updateCmd)

}
