package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/logrusorgru/aurora"
)

type DockerService interface {
	Status() (bool, error)
	Start() ([]byte, error)
	Stop() error
}

type Service struct {
	Name          string
	Address       string
	ContainerName string
	Domain        string
	Shell         string
	ImageName     string
	Cmds          struct {
		RunCmd  []string
		StopCmd []string
		DelCmd  []string
	}
	RunCmd []string

	Config container.Config
	HostConfig container.HostConfig
	NetworkConfig network.NetworkingConfig
}

func DockerRun(ds *Service) ([]byte, error) {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return []byte{}, err
	}

	resp, err := cli.ContainerCreate(ctx, &ds.Config, &ds.HostConfig, &ds.NetworkConfig, ds.ContainerName)
	if err != nil {
		return []byte{}, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return []byte{}, err
	}

	return []byte{}, nil
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

func DockerNetworkConnect(network string, container string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	err = cli.NetworkConnect(ctx, network, container, nil)
	if err != nil {
		return err
	}
	return nil
}

func (ds *Service) Start() ([]byte, error) {

	s, e := ds.Status()
	if e != nil {
		fmt.Println(e)
		return []byte{}, e
	}
	if s {
		if ds.ContainerName != "amazeeio-ssh-agent-add-key" {
			Green(fmt.Sprintf("Already running %v", ds.ContainerName))
		}
		return []byte{}, nil
	}

	container, _ := ds.GetDetails()
	if container.ImageID == "" {
		if !s {
			if ds.ContainerName != "amazeeio-ssh-agent-add-key" {
				_, err := DockerRun(ds)
				if c, _ := ds.GetDetails(); c.ID != "" {
					Green(fmt.Sprintf("Successfully started %v", ds.ContainerName))
				}
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			Red(fmt.Sprintf("Failed to run %v.  Command docker %v failed", ds.ContainerName, strings.Join(ds.RunCmd, " ")))
		}
	}

	return []byte{}, nil
}

func (ds *Service) Status() (bool, error) {

	// amazeeio-ssh-agent-add-key will not show in `docker ps`.
	if ds.ContainerName == "amazeeio-ssh-agent-add-key" {
		return true, nil
	}
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet:   true,
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
		Quiet:   true,
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

	if ds.Name == "" {
		return nil
	}
	names := []string{"/"+ds.Name, ds.Name}

	for _, name := range names {
		if e := DockerKill(name); e == nil {
			Green(fmt.Sprintf("%v container killed", name))
		}
		if e := DockerStop(name); e == nil {
			Green(fmt.Sprintf("%v container stopped", name))
		}
		if e := DockerRemove(name); e != nil {
			Green(fmt.Sprintf("%v container successfully removed", name))
		}
	}

	return nil
}

func (ds *Service) Stop() error {

	container, err := ds.GetDetails()
	if err != nil {
		Green(fmt.Sprintf("Not running %v", ds.ContainerName))
		return nil
	}

	for _, name := range container.Names {
		if e := DockerKill(container.ID); e == nil {
			Green(fmt.Sprintf("%v container killed", name))
		}
		if e := DockerStop(container.ID); e == nil {
			Green(fmt.Sprintf("%v container stopped", name))
		}
		if e := DockerRemove(container.ID); e != nil {
			Green(fmt.Sprintf("%v container successfully removed", name))
		}
	}

	return nil
}

func Red(input string) {
	fmt.Println(aurora.Red(input))
}

func Green(input string) {
	fmt.Println(aurora.Green(input))
}

var _ DockerService = (*Service)(nil)
