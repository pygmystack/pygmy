package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/containerd/containerd/platforms"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/pygmystack/pygmy/service/endpoint"
)

// GetSocket is a bandaid solution to providing a custom socket path.
func GetSocket() string {
	socket := os.Getenv("PYGMY_SOCKET")
	if socket == "" {
		return client.DefaultDockerHost
	}
	if strings.HasPrefix(socket, "unix:///") {
		return socket
	}
	return fmt.Sprintf("unix://%s", socket)
}

// DockerContainerList will return a slice of containers
func DockerContainerList() ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return []types.Container{}, err
	}

	return containers, nil

}

// DockerImageList will return a slice of Docker images.
func DockerImageList() ([]types.ImageSummary, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err)
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return []types.ImageSummary{}, err
	}

	return images, nil

}

// DockerPull will pull a Docker image into the daemon.
func DockerPull(image string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err)
	}

	{

		// To support image references from external sources to docker.io we need to check
		// and validate the image reference for all known cases of validity.

		if m, _ := regexp.MatchString("^([a-zA-Z0-9]+[/][a-zA-Z0-9:-_]+[a-zA-Z0-9:-_.]+)$", image); m {
			// URL was not provided (in full), but the tag was provided.
			// For this, we prepend 'docker.io/' to the reference.
			// Examples:
			//  - amazeeio/pygmy:latest
			image = fmt.Sprintf("docker.io/%v", image)
		} else if m, _ := regexp.MatchString("^([a-zA-Z0-9]+[/][a-zA-Z0-9:-]+[a-zA-Z0-9:-_.]+)$", image); m {
			// URL was not provided (in full), but the tag was not provided.
			// For this, we prepend 'docker.io/' to the reference.
			// Examples:
			//  - amazeeio/pygmy
			image = fmt.Sprintf("docker.io/%v", image)
		} else if m, _ := regexp.MatchString("^([a-zA-Z0-9.].+[a-zA-Z0-9]+[/][a-zA-Z0-9:-_]+[a-zA-Z0-9:-_.]+)$", image); m {
			// URL was provided (in full), but the tag was provided.
			// For this, we do not alter the value provided.
			// Examples:
			//  - quay.io/amazeeio/pygmy:latest
		} else if m, _ := regexp.MatchString("^([a-zA-Z0-9.].+[a-zA-Z0-9]+[/][a-zA-Z0-9:-_]+)$", image); m {
			// URL was provided (in full), but the tag was not provided.
			// For this, we do not alter the value provided.
			// Examples:
			//  - quay.io/amazeeio/pygmy
		} else if m, _ := regexp.MatchString("^([a-zA-Z0-9]+[:][a-zA-Z0-9.-_]+)$", image); m {
			// Library image was provided with tag identifier.
			// For this, we prepend 'docker.io/' to the reference.
			// Examples:
			//  - pygmy:latest
			image = fmt.Sprintf("docker.io/%v", image)
		} else if m, _ := regexp.MatchString("^([a-zA-Z0-9-_]+)$", image); m {
			// Library image was provided without tag identifier.
			// For this, we prepend 'docker.io/' to the reference.
			// Examples:
			//  - pygmy
			image = fmt.Sprintf("docker.io/%v", image)
		} else {
			// Validation not successful
			return "", fmt.Errorf("error: regexp validation for %v failed", image)
		}
	}

	// DockerHub Registry causes a stack trace fatal error when unavailable.
	// We can check for this and report back, handling it gracefully and
	// tell the user the service is down momentarily, and to try again shortly.
	if strings.HasPrefix(image, "docker.io") {
		if s := endpoint.Validate("https://registry-1.docker.io/v2/"); !s {
			return "", fmt.Errorf("cannot reach the Docker Hub Registry, please try again in a few minutes")
		}
	}

	data, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
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
		if strings.Contains(event.Status, fmt.Sprint("Downloaded newer image")) {
			return fmt.Sprintf("Successfully pulled %v", image), nil
		}

		if strings.Contains(event.Status, fmt.Sprint("Image is up to date")) {
			return fmt.Sprintf("Image %v is up to date", image), nil
		}
	}
	return event.Status, nil
}

// DockerStop will stop the container.
func DockerStop(name string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
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

// DockerKill will kill the container.
func DockerKill(name string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	err = cli.ContainerKill(ctx, name, "")
	if err != nil {
		return err
	}
	return nil
}

// DockerRemove will remove the container.
// It will not remove the image.
func DockerRemove(id string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
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
	netVal, _ := DockerNetworkStatus(network.Name)
	if netVal {
		return fmt.Errorf("docker network %v already exists", network.Name)
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
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
		return err
	}

	return nil
}

// DockerNetworkRemove will attempt to remove a Docker network
// and will not apply force to removal.
func DockerNetworkRemove(network string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	err = cli.NetworkRemove(ctx, network)
	if err != nil {
		return err
	}
	return nil
}

// DockerNetworkStatus will identify if a network with a
// specified name is present been created and return a boolean.
func DockerNetworkStatus(network string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}

	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return false, err
	}

	for _, n := range networks {
		if n.Name == network {
			return true, nil
		}
	}

	return false, nil
}

// DockerNetworkGet will use the Docker API to retrieve a Docker network
// which has a given name.
func DockerNetworkGet(name string) (types.NetworkResource, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return types.NetworkResource{}, err
	}
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return types.NetworkResource{}, err
	}
	for _, network := range networks {
		if val, ok := network.Labels["pygmy.name"]; ok {
			if val == name {
				return network, nil
			}
		}
	}
	return types.NetworkResource{}, nil
}

// DockerNetworkConnect will connect a container to a network.
func DockerNetworkConnect(network string, containerName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	e := cli.NetworkConnect(ctx, network, containerName, nil)
	if e != nil {
		return e
	}
	return nil
}

// DockerNetworkConnected will check if a container is connected to a network.
func DockerNetworkConnected(network string, containerName string) (bool, error) {
	// Reset network state:
	c, _ := DockerContainerList()
	for d := range c {
		if c[d].Labels["pygmy.name"] == containerName {
			for net := range c[d].NetworkSettings.Networks {
				if net == network {
					return true, nil
				}
			}
		}
	}
	return false, fmt.Errorf("network was found without the container connected")
}

// DockerVolumeExists will check if a Docker volume has been created.
func DockerVolumeExists(volume types.Volume) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}
	_, _, err = cli.VolumeInspectWithRaw(ctx, volume.Name)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DockerVolumeGet will return the full contents of a types.Volume from the API.
func DockerVolumeGet(name string) (types.Volume, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())

	if err != nil {
		return types.Volume{
			Name: name,
		}, err
	}

	volumes, err := cli.VolumeList(ctx, filters.Args{})
	if err != nil {
		return types.Volume{
			Name: name,
		}, err
	}

	for _, volume := range volumes.Volumes {
		if volume.Name == name {
			return *volume, nil
		}
	}

	return types.Volume{
		Name: name,
	}, nil
}

// DockerVolumeCreate will create a Docker Volume as configured.
func DockerVolumeCreate(volume types.Volume) (types.Volume, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return types.Volume{}, err
	}
	return cli.VolumeCreate(ctx, volumetypes.VolumeCreateBody{
		Driver:     volume.Driver,
		DriverOpts: volume.Options,
		Labels:     volume.Labels,
		Name:       volume.Name,
	})
}

// DockerInspect will return the full container object.
func DockerInspect(container string) (types.ContainerJSON, error) {
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return types.ContainerJSON{}, err
	}

	return cli.ContainerInspect(context.Background(), container)
}

// DockerExec will run a command in a Docker container and return the output.
func DockerExec(container string, command string) ([]byte, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return []byte{}, err
	}

	rst, err := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          strings.Split(command, " ")})

	if err != nil {
		return []byte{}, err
	}

	response, err := cli.ContainerExecAttach(context.Background(), rst.ID, types.ExecStartCheck{})

	if err != nil {
		return []byte{}, err
	}

	data, _ := ioutil.ReadAll(response.Reader)
	defer response.Close()
	return data, nil

}

// DockerContainerCreate will create a container, but will not run it.
func DockerContainerCreate(ID string, config container.Config, hostconfig container.HostConfig, networkconfig network.NetworkingConfig) (container.ContainerCreateCreatedBody, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}
	platform := platforms.Normalize(v1.Platform{
		Architecture: runtime.GOARCH,
		OS:           "linux",
	})
	resp, err := cli.ContainerCreate(ctx, &config, &hostconfig, &networkconfig, &platform, ID)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}
	return resp, err
}

// DockerContainerStart will run an existing container.
func DockerContainerStart(ID string, options types.ContainerStartOptions) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	if err := cli.ContainerStart(ctx, ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return err
}

// DockerContainerLogs will synchronously (blocking, non-concurrently) print
// logs to stdout and stderr, useful for quick containers with a small amount
// of output which are expected to exit quickly.
func DockerContainerLogs(ID string) ([]byte, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost(GetSocket()), client.WithAPIVersionNegotiation())
	if err != nil {
		return []byte{}, err
	}
	b, e := cli.ContainerLogs(ctx, ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})

	if e != nil {
		return []byte{}, e
	}

	buf := new(bytes.Buffer)
	if _, f := buf.ReadFrom(b); f != nil {
		fmt.Println(f)
	}

	return buf.Bytes(), nil
}
