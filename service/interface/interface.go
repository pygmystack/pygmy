package model

import (
	"bytes"
	"fmt"
	"github.com/logrusorgru/aurora"
	"os/exec"
	"strings"
	"sync"
)

type DockerService interface {
	HasDockerClient() bool
	Status() (bool, error)
	Start() ([]byte, error)
	Stop() error
}

type Service struct {
	Name string
	Address string
	ContainerName string
	Domain string
	Shell string
	ImageName string
	Cmds struct {
		RunCmd []string
		StopCmd []string
		DelCmd []string
	}
	RunCmd []string
}

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

func (ds *Service) Start() ([]byte, error) {

	s := true
	if ds.ContainerName != "amazeeio-ssh-agent-key-add" {
		s, e := ds.Status()
		if e != nil {
			return []byte{}, e
		}
		if s {
			Green(fmt.Sprintf("Already running %v", ds.ContainerName))
			return []byte{}, nil
		}
	}

	d, e := DockerRun(ds.RunCmd)

	if s {
		Green(fmt.Sprintf("Successfully started %v", ds.ContainerName))
	} else {
		Red(fmt.Sprintf("Failed to run %v.  Command docker %v failed", ds.ContainerName, strings.Join(ds.RunCmd, " ")))
	}

	return d, e
}

func (ds *Service) Status() (bool, error) {

	data, e := DockerRun([]string{"ps", "--format", "{{.Names}}"})
	if e != nil {
		return false, e
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, ds.ContainerName) {
			return true, nil
		}
	}

	return false, nil

}

func (ds *Service) HasDockerClient() bool { return false }

func (ds *Service) Stop() error {

	s, e := ds.Status()
	if e != nil {
		return e
	}

	if !s {
		Green(fmt.Sprintf("Not running %v", ds.ContainerName))
		return nil
	}

	if _, e := DockerRun([]string{"stop", ds.ContainerName}); e == nil { Green(fmt.Sprintf("%v container stopped", ds.Name)) }
	if _, e := DockerRun([]string{"rm", ds.ContainerName}); e != nil { Green(fmt.Sprintf("%v container successfully deleted", ds.Name)) }

	s, e = ds.PS()
	if e != nil {
		return e
	}

	if !s {
		Green(fmt.Sprintf("Stopped %v", ds.ContainerName))
	} else {
		Red(fmt.Sprintf("Failed to stop %v", ds.ContainerName))
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