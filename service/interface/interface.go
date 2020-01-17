package model

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	volume2 "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type DockerService interface {
	Setup() error
	Status() (bool, error)
	Start() ([]byte, error)
	Stop() error
}

type Service struct {
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
		DockerKill(name)
		DockerRemove(name)
	}

	output, err := DockerRun(Service)
	if err != nil {
		return []byte{}, err
	}

	if c, err := GetRunning(Service); c.ID != "" {
		if !Service.HostConfig.AutoRemove || !discrete {
			fmt.Printf("Successfully started %v\n", name)
		} else if !Service.HostConfig.AutoRemove {
			// We cannot guarantee this container is running at this point if it is to be removed.
			return output, errors.New(fmt.Sprintf("Failed to run %v: %v\n", name, err))
		}
	}

	return output, nil
}

func (Service *Service) Status() (bool, error) {

	name, _ := Service.GetFieldString("name")

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
		for _, n := range container.Names {
			if strings.Contains(n, name) {
				return true, nil
			}
		}
	}

	return false, nil

}

//
func (Service *Service) GetFieldString(field string) (string, error) {

	f := fmt.Sprintf("pygmy.%v", field)

	if container, running := GetRunning(Service); running == nil {
		if val, ok := container.Labels[f]; ok {
			return val, nil
		}
	}

	if val, ok := Service.Config.Labels[f]; ok {
		return val, nil
	}

	return "", errors.New(fmt.Sprintf("could not find field 'pygmy.%v' on service using image %v?", field, Service.Config.Image))
}

func (Service *Service) GetFieldInt(field string) (int, error) {

	f := fmt.Sprintf("pygmy.%v", field)

	if container, running := GetRunning(Service); running == nil {
		if val, ok := container.Labels[f]; ok {
			i, e := strconv.ParseInt(val, 10, 10)
			if e != nil {
				return 0, e
			}
			return int(i), nil
		}
	}

	if val, ok := Service.Config.Labels[f]; ok {
		i, e := strconv.ParseInt(val, 10, 10)
		if e != nil {
			return 0, e
		}
		return int(i), nil
	}

	return 0, errors.New(fmt.Sprintf("could not find field 'pygmy.%v' on service using image %v?", field, Service.Config.Image))
}

func (Service *Service) GetFieldBool(field string) (bool, error) {

	f := fmt.Sprintf("pygmy.%v", field)

	if container, running := GetRunning(Service); running == nil {
		if val, ok := container.Labels[f]; ok {
			if val == "true" {
				return true, nil
			} else if val == "false" {
				return false, nil
			}
		}
	}

	if val, ok := Service.Config.Labels[f]; ok {
		if val == "true" {
			return true, nil
		} else if val == "false" {
			return false, nil
		}
	}

	return false, errors.New(fmt.Sprintf("could not find field 'pygmy.%v' on service using image %v?", field, Service.Config.Image))
}

// GetRunning will get a types.Container variable for a given running container
// and it will not retrieve any information on containers that are not running.
func GetRunning(Service *Service) (types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet: true,
	})

	for _, container := range containers {
		if _, ok := container.Labels["pygmy.name"]; ok {
			if strings.Contains(container.Names[0], Service.Config.Labels["pygmy.name"]) {
				return container, nil
			}
		}
	}
	return types.Container{}, errors.New(fmt.Sprintf("container using image '%v' was not found\n", Service.Config.Image))
}

func (Service *Service) Clean() error {

	pygmy, e := Service.GetFieldString("pygmy")
	name, e := Service.GetFieldString("name")
	if e != nil {
		return nil
	}

	Containers, _ := DockerContainerList()
	for _, container := range Containers {
		if container.Names[0] == name {
			if pygmy == "pygmy" {
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

	name, e := Service.GetFieldString("name")
	discrete, _ := Service.GetFieldBool("discrete")
	if e != nil {
		return nil
	}

	container, err := GetRunning(Service)
	if err != nil {
		if !discrete {
			fmt.Printf("Not running %v\n", name)
		}
		return nil
	}

	for _, name := range container.Names {
		if e := DockerStop(container.ID); e == nil {
			if e := DockerRemove(container.ID); e == nil {
				if !discrete {
					containerName := strings.Trim(name, "/")
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

	// Sanity check to ensure we don't get name conflicts.
	c, _ := DockerContainerList()
	for _, cn := range c {
		if strings.HasSuffix(cn.Names[0], Service.Config.Labels["pygmy.name"])  {
			return []byte{}, nil
		}
	}

	// We need the container name.
	name, e := Service.GetFieldString("name")
	if e != nil {
		return []byte{}, errors.New("container config is missing label for name")
	}

	resp, err := cli.ContainerCreate(ctx, &Service.Config, &Service.HostConfig, &Service.NetworkConfig, name)
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

	b.Close()

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
func DockerNetworkCreate(network *types.NetworkResource) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}

	config := types.NetworkCreate{
		Driver:     network.Driver,
		EnableIPv6: network.EnableIPv6,
		IPAM:       &network.IPAM,
		Internal:   network.Internal,
		Attachable: network.Attachable,
		Options:    network.Options,
		Labels:     network.Labels,
	}

	_, err = cli.NetworkCreate(ctx, network.Name, config)
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

func DockerVolumeExists(volume types.Volume) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return false, err
	}
	_, _, err = cli.VolumeInspectWithRaw(ctx, volume.Name)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DockerVolumeCreate(volume types.Volume) (types.Volume, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return types.Volume{}, err
	}
	return cli.VolumeCreate(ctx, volume2.VolumesCreateBody{
		Driver:     volume.Driver,
		DriverOpts: volume.Options,
		Labels:     volume.Labels,
		Name:       volume.Name,
	})
}
