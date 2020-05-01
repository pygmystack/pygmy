package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	volume2 "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/fubarhouse/pygmy-go/service/endpoint"
)

// DockerService is the requirements for a Docker container to be compatible.
// The Service struct is used to implement this interface, and individual
// variables of type Service can/have overwritten them when logic deems
// it necessary.
type DockerService interface {
	Setup() error
	Status() (bool, error)
	Start() ([]byte, error)
	Stop() error
}

// Service is a collection of requirements for starting a container and
// provides a way for config of any container to be overridden and start
// fully compatible with Docker's API.
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

// Setup will detect if the Service's image reference exists and will
// attempt to run `docker pull` on the non-canonical image if it is
// not found in the daemon.
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

	msg, err := DockerPull(Service.Config.Image)
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
		if e := DockerKill(name); e != nil {
			fmt.Sprintln(e)
		}
		if e := DockerRemove(name); e != nil {
			fmt.Sprintln(e)
		}

	}

	output, err := DockerRun(Service)
	if err != nil {
		return []byte{}, err
	}

	if c, err := GetRunning(Service); c.ID != "" {
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
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	containers, _ := cli.ContainerList(ctx, types.ContainerListOptions{
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

// SetField will set a pygmy label to be equal to the string equal of
// an interface{}, even if it already exists. It should not matter if
// this container is running or not.
func (Service *Service) SetField(name string, value interface{}) error {
	if _, ok := Service.Config.Labels["pygmy."+fmt.Sprint(name)]; !ok {
		//
	} else {
		old, _ := Service.GetFieldString(name)
		Service.Config.Labels["pygmy."+name] = fmt.Sprint(value)
		new, _ := Service.GetFieldString(name)

		if old == new {
			return fmt.Errorf("tag was not set")
		}
	}

	return nil
}

// GetFieldString will get and return a tag on the service using the pygmy
// convention ("pygmy.*") and return it as a string.
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

	return "", fmt.Errorf("could not find field 'pygmy.%v' on service using image %v?", field, Service.Config.Image)
}

// GetFieldInt will get and return a tag on the service using the pygmy
// convention ("pygmy.*") and return it as an int.
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

	return 0, fmt.Errorf("could not find field 'pygmy.%v' on service using image %v?", field, Service.Config.Image)
}

// GetFieldBool will get and return a tag on the service using the pygmy
// convention ("pygmy.*") and return it as a bool.
func (Service *Service) GetFieldBool(field string) (bool, error) {

	f := fmt.Sprintf("pygmy.%v", field)

	if container, running := GetRunning(Service); running == nil {
		if Service.Config.Labels[f] == container.Labels[f] {
			if val, ok := container.Labels[f]; ok {
				if val == "true" {
					return true, nil
				} else if val == "false" {
					return false, nil
				}
			}
		}
	}

	if val, ok := Service.Config.Labels[f]; ok {
		if val == "true" || val == "1" {
			return true, nil
		} else if val == "false" || val == "0" {
			return false, nil
		}
	}

	return false, fmt.Errorf("could not find field 'pygmy.%v' on service using image %v?", field, Service.Config.Image)
}

// GetRunning will get a types.Container variable for a given running container
// and it will not retrieve any information on containers that are not running.
func GetRunning(Service *Service) (types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	containers, _ := cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet: true,
	})

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

	Containers, _ := DockerContainerList()
	for _, container := range Containers {
		if container.Names[0] == name {
			if pygmy {
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

// Stop will stop the container.
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

// _ will ensure DockerService is implemented by Service.
var _ DockerService = (*Service)(nil)

// DockerContainerList will return a slice of containers
func DockerContainerList() ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
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

// DockerPull will pull a Docker image into the daemon.
func DockerPull(image string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
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
		if s := endpoint.Validate("https://registry-1.docker.io"); !s {
			return "", fmt.Errorf("cannot reach the Docker Hub Registry, please try again in a few minutes.")
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
		if strings.Contains(event.Status, fmt.Sprintf("Downloaded newer image for %s", image)) {
			return fmt.Sprintf("Successfully pulled %v\n", image), nil
		}

		if strings.Contains(event.Status, fmt.Sprintf("Image is up to date for %s", image)) {
			return fmt.Sprintf("Image %v is up to date\n", image), nil
		}
	}
	return "", nil
}

// DockerRun will setup and run a given container.
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
		if e := Service.Setup(); e != nil {
			fmt.Println(e)
		}
	}

	// Sanity check to ensure we don't get name conflicts.
	c, _ := DockerContainerList()
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
	if _, f := buf.ReadFrom(b); f != nil {
		fmt.Println(f)
	}

	b.Close()

	return buf.Bytes(), nil
}

// DockerStop will stop the container.
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

// DockerKill will kill the container.
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

// DockerRemove will remove the container.
// It will not remove the image.
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

	if val, ok := network.Labels["pygmy.network"]; ok {
		if network.Name != "" && (val == "true" || val == "1") {
			_, err = cli.NetworkCreate(ctx, network.Name, config)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DockerNetworkRemove will attempt to remove a Docker network
// and will not apply force to removal.
func DockerNetworkRemove(network *types.NetworkResource) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	if val, ok := network.Labels["pygmy.network"]; ok {
		if network.Name != "" && (val == "true" || val == "1") {
			err = cli.NetworkRemove(ctx, network.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DockerNetworkStatus will identify if a network with a
// specified name has been created and return a boolean.
func DockerNetworkStatus(network *types.NetworkResource) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return false, err
	}

	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return false, err
	}

	for _, n := range networks {
		if val, ok := network.Labels["pygmy.network"]; ok {
			if n.Name != "" && (val == "true" || val == "1") {
				return true, nil
			}
		}
	}

	return false, nil
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
		if val, ok := network.Labels["pygmy.network"]; ok {
			if network.Name != "" && (val == "true" || val == "1") {
				return network, nil
			}
		}
	}
	return types.NetworkResource{}, nil
}

// DockerNetworkConnect will connect a container to a network.
func DockerNetworkConnect(network types.NetworkResource, containerName string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	if val, ok := network.Labels["pygmy.network"]; ok {
		if network.Name != "" && (val == "true" || val == "1") {
			e := cli.NetworkConnect(ctx, network.Name, containerName, nil)
			if e != nil {
				fmt.Println(e)
			}
		}
	}
	return nil
}

// DockerNetworkConnect will check if a container is connected to a network.
func DockerNetworkConnected(network types.NetworkResource, containerName string) (bool, error) {
	// Reset network state:
	network, _ = DockerNetworkGet(network.Name)

	for _, container := range network.Containers {
		if container.Name == containerName && container.EndpointID != "" {
			return true, nil
		}
	}
	return false, fmt.Errorf("network was found without the container connected")
}

// DockerVolumeExists will check if a Docker volume has been created.
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

// DockerVolumeGet will return the full contents of a types.Volume from the API.
func DockerVolumeGet(name string) (types.Volume, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()

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

// DockerExec will run a command in a Docker container and return the output.
func DockerExec(container string, command string) ([]byte, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return []byte{}, err
	}

	if rst, err := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          strings.Split(command, " ")}); err != nil {
		return []byte{}, err
	} else {
		if response, err := cli.ContainerExecAttach(context.Background(), rst.ID, types.ExecConfig{}); err != nil {
			return []byte{}, err
		} else {
			data, _ := ioutil.ReadAll(response.Reader)
			defer response.Close()
			return data, nil
		}
	}
}
