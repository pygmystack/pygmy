package model

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"

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
}

//func DockerRun(name string, args []string) ([]byte, error) {
//	ctx := context.Background()
//	cli, err := client.NewEnvClient()
//	if err != nil {
//		fmt.Println(err)
//	}
//	if err := cli.ContainerStart(ctx, name, types.ContainerStartOptions{}); err != nil {
//		fmt.Println(err)
//	}
//	return []byte{}, err
//}

func DockerRun(args []string) ([]byte, error) {

	docker, err := exec.LookPath("docker")
	if err != nil {
		fmt.Println(err)
	}

	// Generate the command, based on input.
	cmd := exec.Cmd{}
	cmd.Path = docker
	cmd.Args = []string{docker}

	// Add our arguments to the command.
	cmd.Args = append(cmd.Args, args...)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	// Check the errors, return as needed.
	var wg sync.WaitGroup
	wg.Add(1)
	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	wg.Done()

	return output.Bytes(), nil
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

func DockerNetworkCreate(name string, args []string) error {
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

func DockerNetworkConnect(source, destination string, args []string) error {
	//ctx := context.Background()
	//cli, err := client.NewEnvClient()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//config := network.EndpointSettings{
	//	EndpointSettings: nil,
	//	IPAMOperational:  false,
	//}
	//err = cli.NetworkConnect(ctx, destination, source, config)
	//if err != nil {
	//	fmt.Println(err)
	//}
	return nil
}

//func DockerRun(args []string) ([]byte, error) {
//
//	docker, err := exec.LookPath("docker")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	// Generate the command, based on input.
//	cmd := exec.Cmd{}
//	cmd.Path = docker
//	cmd.Args = []string{docker}
//
//	// Add our arguments to the command.
//	cmd.Args = append(cmd.Args, args...)
//
//	var output bytes.Buffer
//	cmd.Stdout = &output
//	cmd.Stderr = &output
//
//	// Check the errors, return as needed.
//	var wg sync.WaitGroup
//	wg.Add(1)
//	err = cmd.Run()
//
//	if err != nil {
//		fmt.Println(err)
//		return []byte{}, err
//	}
//	wg.Done()
//
//	return output.Bytes(), nil
//}

func (ds *Service) Start() ([]byte, error) {

	s, e := ds.Status()
	if e != nil {
		return []byte{}, e
	}
	if s {
		if ds.ContainerName != "amazeeio-ssh-agent-add-key" {
			Green(fmt.Sprintf("Already running %v", ds.ContainerName))
		}
		return []byte{}, nil
	}

	d, e := DockerRun(ds.RunCmd)

	if s {
		if ds.ContainerName != "amazeeio-ssh-agent-add-key" {
			Green(fmt.Sprintf("Successfully started %v", ds.ContainerName))
		}
	} else {
		Red(fmt.Sprintf("Failed to run %v.  Command docker %v failed", ds.ContainerName, strings.Join(ds.RunCmd, " ")))
	}

	return d, e
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

func (ds *Service) Stop() error {

	s, e := ds.Status()
	if e != nil {
		return e
	}

	if !s {
		Green(fmt.Sprintf("Not running %v", ds.ContainerName))
		return nil
	}

	if e := DockerKill(ds.ContainerName); e == nil {
		Green(fmt.Sprintf("%v container stopped", ds.Name))
	}
	if e := DockerRemove(ds.ContainerName); e != nil {
		Green(fmt.Sprintf("%v container successfully deleted", ds.Name))
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
