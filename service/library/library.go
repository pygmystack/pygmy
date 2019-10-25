// Library is a package which exposes the commands externally to the compiled binaries.
package library

import (
	"fmt"
	"os"

	"github.com/fubarhouse/pygmy/service/amazee"
	"github.com/fubarhouse/pygmy/service/dnsmasq"
	"github.com/fubarhouse/pygmy/service/haproxy"
	"github.com/fubarhouse/pygmy/service/haproxy_connector"
	model "github.com/fubarhouse/pygmy/service/interface"
	"github.com/fubarhouse/pygmy/service/mailhog"
	"github.com/fubarhouse/pygmy/service/network"
	"github.com/fubarhouse/pygmy/service/resolv"
	"github.com/fubarhouse/pygmy/service/ssh_addkey"
	"github.com/fubarhouse/pygmy/service/ssh_agent"
	"github.com/mitchellh/go-homedir"
)

func SshKeyAdd(args []string) {

	// TODO: key needs to be parametered...
	var key string
	if args[0] != "" {
		if _, err := os.Stat(args[0]); err == nil {
			key = args[0]
		} else {
			fmt.Printf("The file path %v does not exist, or is not readable.\n%v\n", args[0], err)
		}
	} else {
		homedir, _ := homedir.Dir()
		fp := fmt.Sprintf("%v%v.ssh%vid_rsa", homedir, string(os.PathSeparator), string(os.PathSeparator))
		if _, err := os.Stat(fp); err == nil {
			key = fp
		} else {
			fmt.Printf("The file path %v does not exist, or is not readable.\n%v\n", args[0], err)
		}
	}

	if key != "" {
		sshKeyAdder := ssh_addkey.NewAdder(key)
		data, _ := sshKeyAdder.Start()
		sshKeyAdder.Clean()
		fmt.Println(string(data))
	}

}

func Clean(args []string) {
	fmt.Sprint(args)

	dnsmasq := dnsmasq.New()
	dnsmasq.Clean()

	haproxy := haproxy.New()
	haproxy.Clean()

	mailhog := mailhog.New()
	mailhog.Clean()

	sshAgent := ssh_agent.New()
	sshAgent.Clean()

	resolv := resolv.New()
	resolv.Clean()
}

func Restart(args []string) {
	fmt.Sprint(args)
	Stop(args)
	Up(args)
}

func Status(args []string) {
	fmt.Sprint(args)

	dnsmasq := dnsmasq.New()
	if s, _ := dnsmasq.Status(); s {
		model.Green(fmt.Sprintf("[*] Dnsmasq: Running as container %v", dnsmasq.ContainerName))
	} else {
		model.Red(fmt.Sprintf("[ ] Dnsmasq is not running"))
	}

	haproxy := haproxy.New()
	if s, _ := haproxy.Status(); s {
		model.Green(fmt.Sprintf("[*] Haproxy: Haproxy as container %v", haproxy.ContainerName))
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
	if s, _ := mailhog.Status(); s {
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
	if s, _ := sshAgent.Status(); s {
		model.Green(fmt.Sprintf("[*] ssh-agent: Running as docker container %v, loaded keys:", sshAgent.ContainerName))
		sshKeyShower := ssh_addkey.NewShower()
		data, _ := sshKeyShower.Start()
		fmt.Println(string(data))
		sshKeyShower.Clean()
	} else {
		model.Red(fmt.Sprintf("[ ] ssh-agent is not running"))
	}
}

func Stop(args []string) {
	fmt.Sprint(args)

	dnsmasq := dnsmasq.New()
	dnsmasq.Stop()

	haproxy := haproxy.New()
	haproxy.Stop()

	mailhog := mailhog.New()
	mailhog.Stop()

	sshAgent := ssh_agent.New()
	sshAgent.Stop()

	resolv := resolv.New()
	resolv.Clean()
}

func Up(args []string) {
	fmt.Sprint(args)

	dnsmasq := dnsmasq.New()
	dnsmasq.Start()

	haproxy := haproxy.New()
	haproxy.Start()

	netStat, _ := network.Status()
	if !netStat {
		network.Create()
	}
	haproxy_connector.Connect()

	mailhog := mailhog.New()
	mailhog.Start()

	sshAgent := ssh_agent.New()
	sshAgent.Start()

	resolv := resolv.New()
	resolv.Configure()

	SshKeyAdd(args)
}

func Update(args []string) {
	fmt.Sprint(args)
	amazee.AmazeeImagePull()
}

func Version(args []string) {
	fmt.Sprint(args)
	fmt.Println("version called")
}
