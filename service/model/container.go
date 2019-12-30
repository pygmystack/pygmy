package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DockerContainerList will gra a list of all the containers which haven't been removed from the daemon
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

// DockerPull will pull an image.
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

	// Event represents the data object the daemon responds with to print its stdout.
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

// DockerRun will run an image with a series of logic. It will check that
// an image exists, it will grab the image if it doesn't exist, create and
// launch the image, grab the logs all through the daemon's API and return
// the data from the daemon with an error. It will also force the creation
// of a 'pygmy' tag in the container started which allows for containers
// to be identified by pygmy at a later time.
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

// DockerStop will Stop a Docker container, without killing or removing.
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

// DockerKill will kill a Docker container
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

// DockerRemove will remove a Docker container.
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
