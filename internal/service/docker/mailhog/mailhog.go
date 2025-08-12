package mailhog

import (
	"fmt"
	"net"
	"runtime"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	"github.com/pygmystack/pygmy/internal/runtime/docker"
)

// New will provide the standard object for the mailhog container.
func New(c *docker.Params, tlsCertPath string) docker.Service {
	serviceSpec := docker.Service{
		Config: container.Config{
			User: "0",
			ExposedPorts: nat.PortSet{
				"80/tcp":   struct{}{},
				"1025/tcp": struct{}{},
				"8025/tcp": struct{}{},
			},
			Env: []string{
				"MH_UI_BIND_ADDR=0.0.0.0:80",
				"MH_API_BIND_ADDR=0.0.0.0:80",
				"AMAZEEIO=AMAZEEIO",
				fmt.Sprintf("AMAZEEIO_URL=mailhog.%s", c.Domain),
			},
			Image: "pygmystack/mailhog",
			Labels: map[string]string{
				"pygmy.defaults": "true",
				"pygmy.enable":   "true",
				"pygmy.name":     "amazeeio-mailhog",
				"pygmy.network":  "amazeeio-network",
				"pygmy.weight":   "15",
			},
		},
		HostConfig: container.HostConfig{
			AutoRemove: false,
			RestartPolicy: container.RestartPolicy{
				Name:              "unless-stopped",
				MaximumRetryCount: 0,
			},
		},
		NetworkConfig: network.NetworkingConfig{},
	}

	if tlsCertPath != "" {
		serviceSpec.Config.Env = append(serviceSpec.Config.Env, "LAGOON_ROUTE=https://mailhog.docker.amazee.io")
		serviceSpec.Config.Labels["pygmy.url"] = fmt.Sprintf("https://mailhog.%s", c.Domain)
	} else {
		serviceSpec.Config.Env = append(serviceSpec.Config.Env, "LAGOON_ROUTE=http://mailhog.docker.amazee.io")
		serviceSpec.Config.Labels["pygmy.url"] = fmt.Sprintf("http://mailhog.%s", c.Domain)
	}

	return serviceSpec
}

// NewDefaultPorts will provide the standard ports used for merging into the
// mailhog config.
func NewDefaultPorts() docker.Service {
	portConfig := docker.Service{
		HostConfig: container.HostConfig{
			PortBindings: nat.PortMap{
				"1025/tcp": []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: "1025",
					},
				},
			},
		},
	}

	if runtime.GOOS == "darwin" {
		randomPort, _ := GetRandomUnusedPort()
		portConfig.HostConfig.PortBindings["80/tcp"] = []nat.PortBinding{
			{
				HostPort: fmt.Sprint(randomPort),
			},
		}
	}

	return portConfig
}

// GetRandomUnusedPort will identify a random port.
func GetRandomUnusedPort() (int, error) {
	// Let the OS pick an available port on localhost
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	defer ln.Close()

	// Extract the assigned port number
	return ln.Addr().(*net.TCPAddr).Port, nil
}
