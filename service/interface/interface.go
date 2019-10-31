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
	"github.com/docker/docker/client"
)

type DockerService interface {
	Setup() error
	Status() (bool, error)
	Start() ([]byte, error)
	Stop() error
}

type Service struct {
	ContainerName string
	Config        container.Config
	HostConfig    container.HostConfig
	NetworkConfig network.NetworkingConfig
}

func Setup(Service *types.Container) error {
	if Service.Image == "" {
		return nil
	}

	images, _ := DockerImageList()
	for _, image := range images {
		if strings.Contains(fmt.Sprint(image.RepoTags), Service.Image) {
			return nil
		}
	}

	err := DockerPull(Service.Image)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func Start(Service *types.Container) ([]byte, error) {

	s, e := Status(Service)
	if e != nil {
		fmt.Println(e)
		return []byte{}, e
	}


	if s && !Service.HostConfig.AutoRemove {
		fmt.Printf("Already running %v\n", Service.Names[0])
		return []byte{}, nil
	}

	if !s || Service.HostConfig.AutoRemove {

		output, err := DockerRun(Service)

		if c, _ := GetDetails(Service); c.ID != "" {
			if !Service.HostConfig.AutoRemove {
				fmt.Printf("Successfully started %v\n", Service.Names[0])
			}
			return output, nil
		}
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("Failed to run %v.\n", Service.Names[0])
	}

	return []byte{}, nil
}

func Status(Service *types.Container) (bool, error) {

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
			if strings.Contains(name, Service.Names[0]) {
				return true, nil
			}
		}
	}

	return false, nil

}

func GetDetails(Service *types.Container) (types.Container, error) {
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
			if strings.Contains(name, Service.Names[0]) {
				return container, nil
			}
		}
	}
	return types.Container{}, errors.New(fmt.Sprintf("container %v was not found\n", Service.Names[0]))
}

func Clean(Service *types.Container) error {

	if Service.Names[0] == "" {
		return nil
	}
	names := []string{"/" + Service.Names[0], Service.Names[0]}

	for _, name := range names {
		if e := DockerKill(name); e == nil {
			if !Service.HostConfig.AutoRemove {
				fmt.Printf("%v container killed\n", name)
			}
		}
		if e := DockerStop(name); e == nil {
			if !Service.HostConfig.AutoRemove {
				fmt.Printf("%v container stopped\n", name)
			}
		}
		if e := DockerRemove(name); e != nil {
			if !Service.HostConfig.AutoRemove {
				fmt.Printf("%v container successfully removed\n", name)
			}
		}
	}

	return nil
}

func Stop(Service *types.Container) error {

	container, err := GetDetails(Service)
	if err != nil {
		fmt.Printf("Not running %v\n", Service.Names[0])
		return nil
	}

	for _, name := range container.Names {
		if e := DockerStop(container.ID); e == nil {
			if e := DockerRemove(container.ID); e == nil {
				fmt.Printf("%v container successfully removed\n", name)
			}
		}
	}

	return nil
}

//var _ DockerService = (*Service)(nil)

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

func DockerRun(Service *types.Container) ([]byte, error) {

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
		if strings.Contains(Service.Image, fmt.Sprint(image.RepoTags)) {

			// We found the image, we don't need to pull it into the registry.
			imageFound = true

		}

	}

	// If we don't have the image available in the registry, pull it in!
	if !imageFound {
		Setup(Service)
	}

	resp, err := cli.ContainerCreate(ctx, &ds.Config, Service.HostConfig, &Service.NetworkConfig, Service.Names[0])
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

