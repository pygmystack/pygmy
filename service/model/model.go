// Model contains a lot of Docker API abstractions which Pygmy uses.
// It's the core package which has all the connection logic which
// is transferred to the daemon.
package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"strings"
)

// DockerService is an interface which has high-level actions a docker container would need to take in its lifecycle.
type DockerService interface {
	Setup() error
	Status() (bool, error)
	Start() ([]byte, error)
	Stop() error
}

// Service contains all the values and types which Pygmy uses to manage its containers.
// Some of which are Docker API types, however most of Pygmy's business-logic can be found
// within this struct.
type Service struct {
	Name     string
	Group    string
	Disabled bool
	Discrete bool
	Output   bool
	Weight   int
	// URL is a variable defined by the service definition for the general knowledge
	// provided to Pygmy users. Optional, it should be a URL with a port and a path
	// where appropriate.
	URL           string
	Config        container.Config
	HostConfig    container.HostConfig
	NetworkConfig network.NetworkingConfig
}

// Setup is a pre-execution task which will look for the Docker image and if needed will
// fetch the image before returning. At present this does not output anything to indicate
// the image is downloading.
func (Service *Service) Setup() error {
	if Service.Config.Image == "" {
		return nil
	}

	images, _ := DockerImageList()
	for _, image := range images {
		if strings.Contains(fmt.Sprint(image.RepoTags), Service.Config.Image) {
			return nil
		}
	}

	err := DockerPull(Service.Config.Image)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// Start will start the container based on business logic. It will attempt to
// clean images with a certain group/tag to prevent namespace collisions and
// will essentially carry out the default actions or the configuration.
func (Service *Service) Start() ([]byte, error) {

	if Service.Name == "" {
		return []byte{}, nil
	}

	s, e := Service.Status()
	if e != nil {
		fmt.Println(e)
		return []byte{}, e
	}

	if s && !Service.HostConfig.AutoRemove && !Service.Discrete {
		fmt.Printf("Already running %v\n", Service.Name)
		return []byte{}, nil
	}

	if Service.Group == "addkeys" || Service.Group == "showkeys" {
		DockerKill(Service.Name)
		DockerRemove(Service.Name)
	}

	if !s || Service.HostConfig.AutoRemove {

		output, err := DockerRun(Service)
		if err != nil {
			fmt.Println(err)
			return []byte{}, err
		}

		if Service.Output {
			fmt.Println(strings.Trim(string(output), "\n"))
		}

		if c, _ := GetDetails(Service); c.ID != "" {
			if !Service.HostConfig.AutoRemove || !Service.Discrete {
				fmt.Printf("Successfully started %v\n", Service.Name)
			}
			return output, nil
		}
		if err != nil {
			return []byte{}, err
		}
	} else {
		fmt.Printf("Failed to run %v.\n", Service.Name)
	}

	return []byte{}, nil
}

// Status will check if the container has been created with the same name
// as configured. It will check if there's a namespace collision.
func (Service *Service) Status() (bool, error) {

	// If the container doesn't persist we should invalidate the status check.
	if Service.HostConfig.AutoRemove {
		return true, nil
	}
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet: true,
	})
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, Service.Name) {
				return true, nil
			}
		}
	}

	return false, nil

}

// GetDetails will return a types.Container{} object which is created from the
// running container if it matches the desired input (*Service).
func GetDetails(Service *Service) (types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet: true,
	})

	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, Service.Name) {
				return container, nil
			}
		}
	}
	return types.Container{}, errors.New(fmt.Sprintf("container %v was not found\n", Service.Name))
}

// Clean will hard-clean a container by using the API to Kill, Stop and Remove
// the container. It's not a clean exit and it will not force action so it's
// still possible the container will need forcible removal but it's generally
// a very efficient way to 'clean' a system from all Pygmy containers.
func (Service *Service) Clean() error {

	if Service.Name == "" {
		return nil
	}

	Containers, _ := DockerContainerList()
	for _, container := range Containers {
		if container.Names[0] == Service.Name {
			if container.Labels["pygmy"] == "pygmy" {
				name := strings.TrimLeft(container.Names[0], "/")
				if e := DockerKill(container.ID); e == nil {
					if !Service.HostConfig.AutoRemove {
						fmt.Printf("Successfully killed %v\n", name)
					}
				}
				if e := DockerStop(container.ID); e == nil {
					if !Service.HostConfig.AutoRemove {
						fmt.Printf("Successfully stopped %v\n", name)
					}
				}
				if e := DockerRemove(container.ID); e != nil {
					if !Service.HostConfig.AutoRemove {
						fmt.Printf("Successfully removed %v\n", name)
					}
				}
			}
		}
	}

	return nil
}

// Stop will Stop and Remove a configured Pygmy container.
func (Service *Service) Stop() error {

	if Service.Name == "" {
		return nil
	}
	container, err := GetDetails(Service)
	if err != nil {
		if !Service.Discrete {
			fmt.Printf("Not running %v\n", Service.Name)
		}
		return nil
	}

	for _, name := range container.Names {
		if e := DockerStop(container.ID); e == nil {
			if e := DockerRemove(container.ID); e == nil {
				if !Service.Discrete {
					containerName := strings.TrimLeft(name, "/")
					fmt.Printf("Successfully removed %v\n", containerName)
				}
			}
		}
	}

	return nil
}

// Force interface DockerService compliance.
var _ DockerService = (*Service)(nil)
