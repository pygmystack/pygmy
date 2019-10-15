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
	"fmt"
	"github.com/fubarhouse/pygmy/service/dnsmasq"
	"github.com/fubarhouse/pygmy/service/haproxy"
	haproxy_connector "github.com/fubarhouse/pygmy/service/haproxy_connector"
	model "github.com/fubarhouse/pygmy/service/interface"
	"github.com/fubarhouse/pygmy/service/mailhog"
	"github.com/fubarhouse/pygmy/service/network"
	"github.com/fubarhouse/pygmy/service/resolv"
	"github.com/fubarhouse/pygmy/service/ssh_addkey"
	"github.com/fubarhouse/pygmy/service/ssh_agent"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Example: "pygmy status",
	Short: "# Report status of the pygmy services",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		dnsmasq := dnsmasq.New()
		if s, _ := dnsmasq.PS(); s {
			model.Green(fmt.Sprintf("[*] Dnsmasq: Running as container %v", dnsmasq.ContainerName))
		} else {
			model.Red(fmt.Sprintf("[ ] Dnsmasq is not running"))
		}

		haproxy := haproxy.New()
		if s, _ := haproxy.PS(); s {
			model.Green(fmt.Sprintf("[*] Haproxy: Haproxy as container %v", dnsmasq.ContainerName))
		} else {
			model.Red(fmt.Sprintf("[ ] Haproxy is not running"))
		}

		netStat, _ := network.Status()
		if netStat {
			model.Green(fmt.Sprintf("[*] Network: Exists as name amazeeio-network"))
		} else {
			model.Red(fmt.Sprintf("[ ] Network: amazeeio-network does not exist"))
		}

		haproxyStatus, _ := haproxy_connector.Connected()
		if haproxyStatus {
			model.Green(fmt.Sprintf("[*] Network: Haproxy amazeeio-haproxy connected to amazeeio-network"))
		} else {
			model.Red(fmt.Sprintf("[ ] Network: Haproxy amazeeio-haproxy is not connected to amazeeio-network"))
		}

		mailhog := mailhog.New()
		if s, _ := mailhog.PS(); s {
			model.Green(fmt.Sprintf("[*] Mailhog: Running as docker container %v", mailhog.ContainerName))
		} else {
			model.Red(fmt.Sprintf("[ ] Mailhog is not running"))
		}

		resolver := resolv.New()
		if resolver.Status() {
			model.Green(fmt.Sprintf("[*] Resolv is property connected"))
		} else {
			model.Red(fmt.Sprintf("[ ] Resolv is not properly connected"))
		}

		sshAgent := ssh_agent.New()
		if s, _ := sshAgent.PS(); s {
			model.Green(fmt.Sprintf("[*] ssh-agent: Running as docker container %v, loaded keys:", sshAgent.ContainerName))
			sshKeyShower := ssh_addkey.NewShower()
			data, _ := sshKeyShower.Start()
			fmt.Println(string(data))
		} else {
			model.Red(fmt.Sprintf("[ ] ssh-agent is not running"))
		}

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
