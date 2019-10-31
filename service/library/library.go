// Library is a package which exposes the commands externally to the compiled binaries.
package library

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/imdario/mergo"
	"os"

	"github.com/fubarhouse/pygmy/v1/service/amazee"
	"github.com/fubarhouse/pygmy/v1/service/haproxy_connector"
	model "github.com/fubarhouse/pygmy/v1/service/interface"
	"github.com/fubarhouse/pygmy/v1/service/network"
	"github.com/fubarhouse/pygmy/v1/service/resolv"
	"github.com/fubarhouse/pygmy/v1/service/ssh/agent"
	"github.com/spf13/viper"
)

// Config is a struct of configurable options which can
// be passed to package library to configure logic for
// continued abstraction.
type Config struct {
	// Key is the path to the Key which should be added.
	Key string `yaml:"Key"`

	// SkipKey indicates key adding should be skipped.
	SkipKey bool

	// SkipResolver indicates the resolver adding/removal
	// should be skipped - for more specific or manual
	// environment implementations.
	SkipResolver bool `yaml:"DisableResolver"`

	// Services is a []types.Container for an index of all Services.
	Services []types.Container `yaml:"services"`

	// Networks is for network configuration
	Networks []struct{
		Name string `yaml:"name"`
		Containers []string `yaml:"containers"`
	} `yaml:"networks"`

	// Resolvers is for all resolvers
	Resolvers []struct{
		Path string `yaml:"path"`
		Data string `yaml:"contents"`
		// command afterwards?
	} `yaml:"resolvers"`
}

func mergeService(destination types.Container, src *types.Container) (*types.Container, error) {
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
	return &destination, nil
}

func getService(s types.Container, c types.Container) types.Container {
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

	if !agent.Search(key) {

		if key != "" {
			c.Key = key
		}

		data, _ := model.Start(Service)
		fmt.Println(string(data))
	} else {
		fmt.Printf("Already added key file %v.\n", c.Key)
	}
}

func Clean(c Config) {

	Setup(&c)
	Containers, _ := model.DockerContainerList()

	for _, Container := range Containers {
		model.Clean(&Container)
	}

	resolv.New().Clean()
}

func Restart(c Config) {
	Down(c)
	Up(c)
}

func Status(c Config) {

	Setup(&c)

	for _, Service := range c.Services {
		if s, _ := model.Status(&Service); s {
			fmt.Printf("[*] %v: Running as container %v\n", Service.Names[0], Service.Names[0])
		} else {
			fmt.Printf("[ ] %v is not running\n", Service.Names[0])
		}
	}

	if resolv.New().Status() {
		fmt.Printf("[*] Resolv is properly conneted\n")
	} else {
		fmt.Printf("[ ] Resolv is not properly connected\n")
	}

}

func Down(c Config) {

	Setup(&c)

	for _, Service := range c.Services {
		model.Stop(&Service)
	}

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

	for _, Service := range c.Services {
		model.Start(&Service)
	}

	for _, Network := range c.Networks {
		netStat, _ := network.Status(Network.Name)
		if !netStat {
			network.Create(Network.Name)
		}
		for _, Container := range Network.Containers {
			if s, _ := haproxy_connector.Connected(Container, Network.Name); !s {
				fmt.Printf("Connecting %v to %v\n", Container, Network.Name)
				haproxy_connector.Connect(Container, Network.Name)
			}
		}
	}

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
