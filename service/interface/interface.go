package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/fubarhouse/pygmy-go/service/interface/docker"
)



// Setup will detect if the Service's image reference exists and will
// attempt to run `docker pull` on the non-canonical image if it is
// not found in the daemon.
func (Service *Service) Setup() error {
	if Service.Config.Image == "" {
		return nil
	}

	images, _ := docker.DockerImageList()
	for _, image := range images {
		if strings.Contains(fmt.Sprint(image.RepoTags), Service.Config.Image) {
			return nil
		}
	}

	msg, err := docker.DockerPull(Service.Config.Image)
	if msg != "" {
		fmt.Println(msg)
	}

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// Start will perform a series of checks to see if the container starting
// is supposed be removed before-hand and will check to see if the
// container is running before it is actually started.
func (Service *Service) Start() ([]byte, error) {

	name, err := Service.GetFieldString("name")
	discrete, _ := Service.GetFieldBool("discrete")
	purpose, _ := Service.GetFieldString("purpose")

	if err != nil {
		return []byte{}, nil
	}

	s := false

	if !Service.HostConfig.AutoRemove {
		var e error
		s, e = Service.Status()
		if e != nil {
			return []byte{}, e
		}
	}

	if s && !Service.HostConfig.AutoRemove && !discrete {
		fmt.Printf("Already running %v\n", name)
		return []byte{}, nil
	}

	if purpose == "addkeys" || purpose == "showkeys" {
		if e := docker.DockerKill(name); e != nil {
			fmt.Sprintln(e)
		}
		if e := docker.DockerRemove(name); e != nil {
			fmt.Sprintln(e)
		}

	}

	output, err := Service.DockerRun()
	if err != nil {
		return []byte{}, err
	}

	if c, err := Service.GetRunning(); c.ID != "" {
		if !Service.HostConfig.AutoRemove && !discrete {
			fmt.Printf("Successfully started %v\n", name)
		} else if Service.HostConfig.AutoRemove && err != nil {
			// We cannot guarantee this container is running at this point if it is to be removed.
			return output, fmt.Errorf("Failed to run %v: %v\n", name, err)
		}
	}

	return output, nil
}

// Status will check if the container is running.
func (Service *Service) Status() (bool, error) {

	name, _ := Service.GetFieldString("name")

	// If the container doesn't persist we should invalidate the status check.
	// This assumes state of any containr with status checks to pass if they
	// are configured with HostConfig.AutoRemove
	if Service.HostConfig.AutoRemove {
		return true, nil
	}
	containers, _ := docker.DockerContainerList()
	for _, container := range containers {
		for _, n := range container.Names {
			if strings.Contains(n, name) {
				return true, nil
			}
		}
	}

	return false, nil

}

// GetRunning will get a types.Container variable for a given running container
// and it will not retrieve any information on containers that are not running.
func (Service *Service) GetRunning() (types.Container, error) {
	containers, _ := docker.DockerContainerList()
	for _, container := range containers {
		if _, ok := container.Labels["pygmy.name"]; ok {
			if strings.Contains(container.Names[0], Service.Config.Labels["pygmy.name"]) {
				return container, nil
			}
		}
	}
	return types.Container{}, fmt.Errorf("container using image '%v' was not found\n", Service.Config.Image)
}

// Clean will cleanup and remove the container.
func (Service *Service) Clean() error {

	pygmy, _ := Service.GetFieldBool("pygmy.enable")
	name, e := Service.GetFieldString("name")
	if e != nil {
		return nil
	}

	Containers, _ := docker.DockerContainerList()
	for _, container := range Containers {
		if container.Names[0] == name {
			if pygmy {
				name := strings.TrimLeft(container.Names[0], "/")
				if e := docker.DockerKill(container.ID); e == nil {
					if !Service.HostConfig.AutoRemove {
						fmt.Printf("Successfully killed %v\n", name)
					}
				}
				if e := docker.DockerStop(container.ID); e == nil {
					if !Service.HostConfig.AutoRemove {
						fmt.Printf("Successfully stopped %v\n", name)
					}
				}
				if e := docker.DockerRemove(container.ID); e != nil {
					if !Service.HostConfig.AutoRemove {
						fmt.Printf("Successfully removed %v\n", name)
					}
				}
			}
		}
	}

	return nil
}

// Stop will stop the container.
func (Service *Service) Stop() error {

	name, e := Service.GetFieldString("name")
	discrete, _ := Service.GetFieldBool("discrete")
	if e != nil {
		return nil
	}

	container, err := Service.GetRunning()
	if err != nil {
		if !discrete {
			fmt.Printf("Not running %v\n", name)
		}
		return nil
	}

	for _, name := range container.Names {
		if e := docker.DockerStop(container.ID); e == nil {
			if e := docker.DockerRemove(container.ID); e == nil {
				if !discrete {
					containerName := strings.Trim(name, "/")
					fmt.Printf("Successfully removed %v\n", containerName)
				}
			}
		}
	}

	return nil
}

// _ will ensure DockerService is implemented by Service.
var _ DockerService = (*Service)(nil)

// DockerRun will setup and run a given container.
func (Service *Service) DockerRun() ([]byte, error) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	cli.NegotiateAPIVersion(ctx)
	if err != nil {
		return []byte{}, err
	}

	// Ensure we have the image available:
	images, _ := docker.DockerImageList()

	// Specify a false boolean which we can switch to true if the image is in the registry:
	imageFound := false

	// Loop over our images
	for _, image := range images {

		// Check if it contains the desired string
		if strings.Contains(Service.Config.Image, fmt.Sprint(image.RepoTags)) {

			// We found the image, we don't need to pull it into the registry.
			imageFound = true

		}

	}

	// If we don't have the image available in the registry, pull it in!
	if !imageFound {
		if e := Service.Setup(); e != nil {
			fmt.Println(e)
		}
	}

	// Sanity check to ensure we don't get name conflicts.
	c, _ := docker.DockerContainerList()
	for _, cn := range c {
		if strings.HasSuffix(cn.Names[0], Service.Config.Labels["pygmy.name"]) {
			return []byte{}, nil
		}
	}

	// We need the container name.
	name, e := Service.GetFieldString("name")
	if e != nil {
		return []byte{}, fmt.Errorf("container config is missing label for name")
	}

	resp, err := docker.DockerContainerCreate(name, Service.Config, Service.HostConfig, Service.NetworkConfig)
	if err != nil {
		return []byte{}, err
	}

	return docker.DockerContainerLogs(resp.ID)

}

