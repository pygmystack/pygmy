// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"github.com/fubarhouse/pygmy/service/dnsmasq"
	"github.com/fubarhouse/pygmy/service/haproxy"
	"github.com/fubarhouse/pygmy/service/mailhog"
	"github.com/fubarhouse/pygmy/service/resolv"
	"github.com/fubarhouse/pygmy/service/ssh_addkey"
	"github.com/fubarhouse/pygmy/service/ssh_agent"
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Example: "pygmy restart",
	Short: "# Report status of the pygmy services",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		dnsmasq := dnsmasq.New()
		haproxy := haproxy.New()
		mailhog := mailhog.New()
		resolver := resolv.New()
		sshAgent := ssh_agent.New()
		sshKeyAdder := ssh_addkey.NewAdder("")

		dnsmasq.Stop()
		haproxy.Stop()
		mailhog.Stop()
		resolver.Clean()
		sshAgent.Stop()

		dnsmasq.Start()
		haproxy.Start()
		mailhog.Start()
		resolver.Configure()
		sshAgent.Start()
		sshKeyAdder.Start()

	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
