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
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
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
		c.SshKeyAdder = getService(ssh_addkey.NewAdder(key), c.SshKeyAdder)

		if key != "" {
			c.Key = key
		}

		data, _ := c.SshKeyAdder.Start()
		fmt.Println(string(data))
	} else {
		fmt.Printf("Already added key file %v.\n", c.Key)
	}
}

func Clean(c Config) {

	Setup(&c)

	c.DnsMasq = getService(dnsmasq.New(), c.DnsMasq)
	c.DnsMasq.Clean()

	c.HaProxy = getService(haproxy.New(), c.HaProxy)
	c.HaProxy.Clean()

	c.MailHog = getService(mailhog.New(), c.MailHog)
	c.MailHog.Clean()

	c.SshAgent = getService(ssh_agent.New(), c.SshAgent)
	c.SshAgent.Clean()

	resolv.New().Clean()
}

func Restart(c Config) {
	Down(c)
	Up(c)
}

func Status(c Config) {

	Setup(&c)

	c.DnsMasq = getService(dnsmasq.New(), c.DnsMasq)
	if s, _ := c.DnsMasq.Status(); s {
		fmt.Printf("[*] Dnsmasq: Running as container %v\n", c.DnsMasq.ContainerName)
	} else {
		fmt.Printf("[ ] Dnsmasq is not running\n")
	}

	c.HaProxy = getService(haproxy.New(), c.HaProxy)
	if s, _ := c.HaProxy.Status(); s {
		fmt.Printf("[*] Haproxy: Haproxy as container %v\n", c.HaProxy.ContainerName)
	} else {
		fmt.Printf("[ ] Haproxy is not running")
	}

	if s, _ := network.Status(c.Network); s {
		fmt.Printf("[*] Network: Exists as name %v\n", c.Network)
	} else {
		fmt.Printf("[ ] Network: %v does not exist\n", c.Network)
	}

	if s, _ := haproxy_connector.Connected(c.HaProxy.ContainerName, c.Network); s {
		fmt.Printf("[*] Network: Haproxy %v connected to %v\n", c.HaProxy.ContainerName, c.Network)
	} else {
		fmt.Printf("[ ] Network: Haproxy %v is not connected to %v\n", c.HaProxy.ContainerName, c.Network)
	}

	c.MailHog = getService(mailhog.New(), c.MailHog)
	if s, _ :=  c.MailHog.Status(); s {
		fmt.Printf("[*] Mailhog: Running as docker container %v\n", c.MailHog.ContainerName)
	} else {
		fmt.Printf("[ ] Mailhog is not running\n")
	}

	if resolv.New().Status() {
		fmt.Printf("[*] Resolv is properly conneted\n")
	} else {
		fmt.Printf("[ ] Resolv is not properly connected\n")
	}

	c.SshAgent = getService(ssh_agent.New(), c.SshAgent)
	if s, _ := c.SshAgent.Status(); s {
		fmt.Printf("[*] ssh-agent: Running as docker container %v, loaded keys:\n", c.SshAgent.ContainerName)
		c.SshKeyLister = getService(ssh_addkey.NewShower(), c.SshKeyLister)
		data, _ := c.SshKeyLister.Start()
		fmt.Println(string(data))
		c.SshKeyLister.Clean()
	} else {
		fmt.Printf("[ ] ssh-agent is not running\n")
	}
}

func Down(c Config) {

	Setup(&c)
	
	c.DnsMasq = getService(dnsmasq.New(), c.DnsMasq)
	c.DnsMasq.Stop()

	c.HaProxy = getService(haproxy.New(), c.HaProxy)
	c.HaProxy.Stop()

	c.MailHog = getService(mailhog.New(), c.MailHog)
	c.MailHog.Stop()

	c.SshAgent = getService(ssh_agent.New(), c.SshAgent)
	c.SshAgent.Stop()

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

	c.DnsMasq = getService(dnsmasq.New(), c.DnsMasq)
	c.DnsMasq.Start()

	c.HaProxy = getService(haproxy.New(), c.HaProxy)
	c.HaProxy.Start()

	netStat, _ := network.Status(c.Network)
	if !netStat {
		network.Create(c.Network)
	}

	if s, _ := haproxy_connector.Connected(c.HaProxy.ContainerName, c.Network); !s {
		fmt.Printf("Connecting %v to %v\n", c.HaProxy.ContainerName, c.Network)
		haproxy_connector.Connect(c.HaProxy.ContainerName, c.Network)
	}

	c.MailHog = getService(mailhog.New(), c.MailHog)
	c.MailHog.Start()

	c.SshAgent = getService(ssh_agent.New(), c.SshAgent)
	c.SshAgent.Start()

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
