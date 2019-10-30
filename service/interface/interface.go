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

func (ds *Service) Setup() error {
	if ds.Config.Image == "" {
		return nil
	}

	images, _ := DockerImageList()
	for _, image := range images {
		if strings.Contains(fmt.Sprint(image.RepoTags), ds.Config.Image) {
			return nil
		}
	}

	err := DockerPull(ds.Config.Image)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (ds *Service) Start() ([]byte, error) {

	s, e := ds.Status()
	if e != nil {
		fmt.Println(e)
		return []byte{}, e
	}


	if s && !ds.HostConfig.AutoRemove {
		fmt.Printf("Already running %v\n", ds.ContainerName)
		return []byte{}, nil
	}

	if !s || ds.HostConfig.AutoRemove {

		output, err := DockerRun(ds)

		if c, _ := ds.GetDetails(); c.ID != "" {
			if !ds.HostConfig.AutoRemove {
				fmt.Printf("Successfully started %v\n", ds.ContainerName)
			}
			return output, nil
		}
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("Failed to run %v.\n", ds.ContainerName)
	}

	return []byte{}, nil
}

func (ds *Service) Status() (bool, error) {

	// If the container doesn't persist we should invalidate the status check.
	if ds.HostConfig.AutoRemove {
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
			if strings.Contains(name, ds.ContainerName) {
				return true, nil
			}
		}
	}

	return false, nil

}

func (ds *Service) GetDetails() (types.Container, error) {
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
			if strings.Contains(name, ds.ContainerName) {
				return container, nil
			}
		}
	}
	return types.Container{}, errors.New(fmt.Sprintf("container %v was not found\n", ds.ContainerName))
}

func (ds *Service) Clean() error {

	if ds.ContainerName == "" {
		return nil
	}
	names := []string{"/" + ds.ContainerName, ds.ContainerName}

	for _, name := range names {
		if e := DockerKill(name); e == nil {
			if !ds.HostConfig.AutoRemove {
				fmt.Printf("%v container killed\n", name)
			}
		}
		if e := DockerStop(name); e == nil {
			if !ds.HostConfig.AutoRemove {
				fmt.Printf("%v container stopped\n", name)
			}
		}
		if e := DockerRemove(name); e != nil {
			if !ds.HostConfig.AutoRemove {
				fmt.Printf("%v container successfully removed\n", name)
			}
		}
	}

	return nil
}

func (ds *Service) Stop() error {

	container, err := ds.GetDetails()
	if err != nil {
		fmt.Printf("Not running %v\n", ds.ContainerName)
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

var _ DockerService = (*Service)(nil)

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

func DockerRun(ds *Service) ([]byte, error) {

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
		if strings.Contains(ds.Config.Image, fmt.Sprint(image.RepoTags)) {

			// We found the image, we don't need to pull it into the registry.
			imageFound = true

		}

	}

	// If we don't have the image available in the registry, pull it in!
	if !imageFound {
		ds.Setup()
	}

	resp, err := cli.ContainerCreate(ctx, &ds.Config, &ds.HostConfig, &ds.NetworkConfig, ds.ContainerName)
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

