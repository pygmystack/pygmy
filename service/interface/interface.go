package model

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type DockerService interface {
	Setup() error
	Status() (bool, error)
	Start() ([]byte, error)
	Stop() error
}

type Service struct {
	Name     string
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

// Network is a struct containing the configuration of a single Docker network
// including some extra fields so that Pygmy knows how to interact with the
// desired outcome.
type Network struct {
	// Name is the name of the network, it is independent of the map key which
	// will be used to configure pygmy but this field should match the map key.
	Name string `yaml:"name"`
	// Containers is a []string which indicates the names of the containers
	// that need to be connected to this network.
	Containers []string `yaml:"containers"`
	// Config is the actual Network configuration for the Docker Network.
	// It is the Network creation configuration as provided by the Docker API.
	Config types.NetworkCreate `yaml:"config"`
}

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

	vOne, _ := Service.TagGet("addkeys")
	vTwo, _ := Service.TagGet("showkeys")
	if vOne == "pygmy.addkeys" || vTwo == "pygmy.showkeys" {
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

// TagGet will return the value and error of obtaining a tag with a given name
// from a container. Pygmy uses tags of a dot notation following the prefix of
// "pygmy", so this will search for "pygmy" when the name parameter is empty,
// or when a value is provided it will look for the tag "pygmy.name" which is
// said to exist on the current container configuration, returning an error if
// the specified tag is not found, otherwise it will also return the value of
// the given tag.
func (Service *Service) TagGet(name string) (string, error) {
	c, _ := DockerContainerList()
	var searchString string
	if name == "" {
		searchString = "pygmy"
	} else {
		searchString = fmt.Sprintf("pygmy.%v", name)
	}
	for _, x := range c {
		for _, n := range x.Names {
			if n == Service.Name {
				for label, value := range x.Labels {
					if label == searchString {
						return value, nil
					}
				}
			}
		}
	}
	return "", errors.New(fmt.Sprintf("container %v does not have the tag %v", Service.Name, name))
}

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

var _ DockerService = (*Service)(nil)

func DockerContainerList() ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return []types.Container{}, err
	}

	return containers, nil

}

func DockerImageList() ([]types.ImageSummary, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return []types.ImageSummary{}, err
	}

	return images, nil

}

func DockerPull(image string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}

	data, err := cli.ImagePull(ctx, "docker.io/"+image, types.ImagePullOptions{})
	if err != nil {
		fmt.Println(err)
	}

	d := json.NewDecoder(data)

	type Event struct {
		Status         string `json:"status"`
		Error          string `json:"error"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"progressDetail"`
	}

	var event *Event
	for {
		if err := d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}
	}

	if event != nil {
		if strings.Contains(event.Status, fmt.Sprintf("Downloaded newer image for %s", image)) {
			fmt.Printf("Successfully pulled %v\n", image)
		}

		if strings.Contains(event.Status, fmt.Sprintf("Image is up to date for %s", image)) {
			fmt.Printf("Image %v is up to date\n", image)
		}
	}
	return nil
}

func DockerRun(Service *Service) ([]byte, error) {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return []byte{}, err
	}

	// Ensure we have the image available:
	images, _ := DockerImageList()

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
		Service.Setup()
	}

	// All pygmy services need some sort of reference for pygmy to consume:
	if Service.Config.Labels["pygmy"] != "pygmy" {
		if Service.Config.Labels == nil {
			Service.Config.Labels = make(map[string]string)
		}
		Service.Config.Labels["pygmy"] = "pygmy"
	}

	resp, err := cli.ContainerCreate(ctx, &Service.Config, &Service.HostConfig, &Service.NetworkConfig, Service.Name)
	if err != nil {
		return []byte{}, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return []byte{}, err
	}

	b, _ := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})

	buf := new(bytes.Buffer)
	buf.ReadFrom(b)

	return buf.Bytes(), nil
}

func DockerStop(name string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	timeout := time.Duration(10)
	err = cli.ContainerStop(ctx, name, &timeout)
	if err != nil {
		return err
	}
	return nil
}

func DockerKill(name string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	err = cli.ContainerKill(ctx, name, "")
	if err != nil {
		return err
	}
	return nil
}

func DockerRemove(id string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	err = cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

// DockerNetworkCreate is an abstraction layer on top of the Docker API call
// which will create a Docker network using a specified configuration.
func DockerNetworkCreate(name string, config types.NetworkCreate) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	_, err = cli.NetworkCreate(ctx, name, config)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// DockerNetworkGet will use the Docker API to retrieve a Docker network
// which has a given name.
func DockerNetworkGet(name string) (types.NetworkResource, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return types.NetworkResource{}, err
	}
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return types.NetworkResource{}, err
	}
	for _, network := range networks {
		if network.Name == name {
			return network, nil
		}
	}
	return types.NetworkResource{}, nil
}

func DockerNetworkConnect(network string, containerName string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	err = cli.NetworkConnect(ctx, network, containerName, nil)
	if err != nil {
		return err
	}
	return nil
}

func DockerVolumeExists(name string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return false, err
	}
	_, _, err = cli.VolumeInspectWithRaw(ctx, name)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DockerVolumeCreate(name string) (types.Volume, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return types.Volume{}, err
	}
	return cli.VolumeCreate(ctx, volume.VolumesCreateBody{
		Name: name,
	})
}
