// Library is a package which exposes the commands externally to the compiled binaries.
package library

import (
	"fmt"
	"github.com/imdario/mergo"
	"os"

	"github.com/fubarhouse/pygmy/service/amazee"
	"github.com/fubarhouse/pygmy/service/dnsmasq"
	"github.com/fubarhouse/pygmy/service/haproxy"
	"github.com/fubarhouse/pygmy/service/haproxy_connector"
	model "github.com/fubarhouse/pygmy/service/interface"
	"github.com/fubarhouse/pygmy/service/mailhog"
	"github.com/fubarhouse/pygmy/service/network"
	"github.com/fubarhouse/pygmy/service/resolv"
	"github.com/fubarhouse/pygmy/service/ssh/ssh_addkey"
	"github.com/fubarhouse/pygmy/service/ssh/ssh_agent"
	"github.com/spf13/viper"
)

// Config is a struct of configurable options which can
// be passed to package library to configure logic for
// continued abstraction.
type Config struct {
	// Key is the path to the Key which should be added.
	Key string `yaml:"Key"`

	// Network
	Network string `yaml:"Network"`

	// SkipKey indicates key adding should be skipped.
	SkipKey bool

	// SkipResolver indicates the resolver adding/removal
	// should be skipped - for more specific or manual
	// environment implementations.
	SkipResolver bool `yaml:"DisableResolver"`

	SshAgent     model.Service `yaml:"SshAgent"`
	DnsMasq      model.Service `yaml:"DnsMasq"`
	HaProxy      model.Service `yaml:"HaProxy"`
	MailHog      model.Service `yaml:"MailHog"`
	SshKeyAdder  model.Service `yaml:"SshKeyAdder"`
	SshKeyLister model.Service `yaml:"SshKeyLister"`
}

func mergeService(destination model.Service, src *model.Service) (*model.Service, error) {
	if err := mergo.Merge(&destination, src); err != nil {
		fmt.Println(err)
		return src, err
	}
	fmt.Println(destination)
	return &destination, nil
}

func getService(s model.Service, c model.Service) model.Service {
	Service, _ := mergeService(s, &c)
	return *Service
}

func SshKeyAdd(c Config, key string) {

	if c.SkipKey {
		return
	}

	Setup(&c)

	if key != "" {
		if _, err := os.Stat(key); err != nil {
			fmt.Printf("The file path %v does not exist, or is not readable.\n%v\n", key, err)
			return
		}
	}

	if !ssh_agent.Search(key) {
		SshKeyService := getService(ssh_addkey.NewAdder(key), c.SshKeyAdder)

		if key != "" {
			c.Key = key
		}

		data, _ := SshKeyService.Start()
		fmt.Println(string(data))
	} else {
		fmt.Printf("Already added key file %v.\n", c.Key)
	}
}

func Clean(c Config) {

	Setup(&c)

	dnsmasq := getService(dnsmasq.New(), c.DnsMasq)
	dnsmasq.Clean()

	haproxy := getService(haproxy.New(), c.HaProxy)
	haproxy.Clean()

	mailhog := getService(mailhog.New(), c.MailHog)
	mailhog.Clean()

	sshagent := getService(ssh_agent.New(), c.SshAgent)
	sshagent.Clean()

	resolv.New().Clean()
}

func Restart(c Config) {
	Down(c)
	Up(c)
}

func Status(c Config) {

	Setup(&c)

	dnsmasq := getService(dnsmasq.New(), c.DnsMasq)
	if s, _ := dnsmasq.Status(); s {
		model.Green(fmt.Sprintf("[*] Dnsmasq: Running as container %v", c.DnsMasq.ContainerName))
	} else {
		model.Red(fmt.Sprintf("[ ] Dnsmasq is not running"))
	}

	haproxy := getService(haproxy.New(), c.HaProxy)
	if s, _ := haproxy.Status(); s {
		model.Green(fmt.Sprintf("[*] Haproxy: Haproxy as container %v", c.HaProxy.ContainerName))
	} else {
		model.Red(fmt.Sprintf("[ ] Haproxy is not running"))
	}

	if s, _ := network.Status(c.Network); s {
		model.Green(fmt.Sprintf("[*] Network: Exists as name %v", c.Network))
	} else {
		model.Red(fmt.Sprintf("[ ] Network: %v does not exist", c.Network))
	}

	if s, _ := haproxy_connector.Connected(c.HaProxy.ContainerName, c.Network); s {
		model.Green(fmt.Sprintf("[*] Network: Haproxy %v connected to %v", c.HaProxy.ContainerName, c.Network))
	} else {
		model.Red(fmt.Sprintf("[ ] Network: Haproxy %v is not connected to %v", c.HaProxy.ContainerName, c.Network))
	}

	mailhog := getService(mailhog.New(), c.MailHog)
	if s, _ :=  mailhog.Status(); s {
		model.Green(fmt.Sprintf("[*] Mailhog: Running as docker container %v", c.MailHog.ContainerName))
	} else {
		model.Red(fmt.Sprintf("[ ] Mailhog is not running"))
	}

	if resolv.New().Status() {
		model.Green(fmt.Sprintf("[*] Resolv is property connected"))
	} else {
		model.Red(fmt.Sprintf("[ ] Resolv is not properly connected"))
	}

	SshAgentService := getService(ssh_agent.New(), c.SshAgent)
	if s, _ := SshAgentService.Status(); s {
		model.Green(fmt.Sprintf("[*] ssh-agent: Running as docker container %v, loaded keys:", SshAgentService.ContainerName))
		sshKeyShower := getService(ssh_addkey.NewShower(), c.SshAgent)
		data, _ := sshKeyShower.Start()
		fmt.Println(string(data))
		sshKeyShower.Clean()
	} else {
		model.Red(fmt.Sprintf("[ ] ssh-agent is not running"))
	}
}

func Down(c Config) {

	Setup(&c)
	
	DnsMasqService := getService(dnsmasq.New(), c.DnsMasq)
	DnsMasqService.Stop()

	HaProxyService := getService(haproxy.New(), c.HaProxy)
	HaProxyService.Stop()

	MailHogService := getService(mailhog.New(), c.MailHog)
	MailHogService.Stop()

	SshAgentService := getService(ssh_agent.New(), c.SshAgent)
	SshAgentService.Stop()

	if !c.SkipResolver {
		resolv := resolv.New()
		resolv.Clean()
	}
}

func Setup(c *Config) {
	viper.SetDefault("Network", "amazeeio-network")
	viper.SetDefault("HaProxy.HostConfig.PortBindings", "map[80/tcp:[map[HostPort:80]]]")

	e := viper.Unmarshal(&c)

	if e != nil {
		fmt.Println(e)
	}

}

func Up(c Config) {

	Setup(&c)

	DnsMasqService := getService(dnsmasq.New(), c.DnsMasq)
	DnsMasqService.Start()

	HaProxyService := getService(haproxy.New(), c.HaProxy)
	HaProxyService.Start()

	netStat, _ := network.Status(c.Network)
	if !netStat {
		network.Create(c.Network)
	}

	if s, _ := haproxy_connector.Connected(c.HaProxy.ContainerName, c.Network); !s {
		haproxy_connector.Connect(c.HaProxy.ContainerName, c.Network)
	}

	MailHogService := getService(mailhog.New(), c.MailHog)
	MailHogService.Start()

	SshAgentService := getService(ssh_agent.New(), c.SshAgent)
	SshAgentService.Start()

	if !c.SkipResolver {
		resolv := resolv.New()
		resolv.Configure()
	}

	if !c.SkipKey {

		SshKeyAdd(c, c.Key)
	}
}

func Update(c Config) {
	amazee.AmazeeImagePull()
}

func Version(c Config) {
	fmt.Println("version called")
}
