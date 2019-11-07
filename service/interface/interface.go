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
	Name          string
	Disabled      bool
	Discrete      bool
	Output        bool
	Config        container.Config
	HostConfig    container.HostConfig
	NetworkConfig network.NetworkingConfig
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

	s, e := Service.Status()
	if e != nil {
		fmt.Println(e)
		return []byte{}, e
	}

	if s && !Service.HostConfig.AutoRemove && !Service.Discrete {
		fmt.Printf("Already running %v\n", Service.Name)
		return []byte{}, nil
	}

	if !s || Service.HostConfig.AutoRemove {

		output, err := DockerRun(Service)

		if Service.Output {
			fmt.Println(string(output))
		}

		if c, _ := GetDetails(Service); c.ID != "" {
			if !Service.HostConfig.AutoRemove || !Service.Discrete {
				fmt.Printf("Successfully started %v\n", Service.Name)
			}
			return output, nil
		}
		if err != nil {
			fmt.Println(err)
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
	names := []string{"/" + Service.Name, Service.Name}

	for _, name := range names {
		if e := DockerKill(name); e == nil {
			if !Service.HostConfig.AutoRemove {
				fmt.Printf("Successfully killed %v\n", name)
			}
		}
		if e := DockerStop(name); e == nil {
			if !Service.HostConfig.AutoRemove {
				fmt.Printf("Successfully stopped %v\n", name)
			}
		}
		if e := DockerRemove(name); e != nil {
			if !Service.HostConfig.AutoRemove {
				fmt.Printf("Successfully removed %v\n", name)
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
					fmt.Printf("Successfully removed %v\n", name)
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

func DockerPull(image string) (error) {
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
	//if Service.Config.Labels["pygmy"] != "pygmy" {
	//	Service.Config.Labels["pygmy"] = "pygmy"
	//}

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

func DockerRemove(name string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	err = cli.ContainerRemove(ctx, name, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

func DockerNetworkCreate(name string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	_, err = cli.NetworkCreate(ctx, name, types.NetworkCreate{})
	if err != nil {
		fmt.Println(err)
	}
	return nil
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
		Name:       name,
	})
}