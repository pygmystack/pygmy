package docker

import (
	"context"
	"fmt"
	"github.com/pygmystack/pygmy/internal/services/color"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	. "github.com/logrusorgru/aurora"
	"golang.org/x/term"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/containers"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/images"
)

// Setup will detect if the Service's image reference exists and will
// attempt to run `docker pull` on the non-canonical image if it is
// not found in the daemon.
func (Service *Service) Setup() error {
	if Service.Config.Image == "" {
		return fmt.Errorf("image reference is nil value")
	}

	found := false
	imgs, _ := images.List()
	for _, image := range imgs {
		if strings.Contains(fmt.Sprint(image.RepoTags), Service.Config.Image) {
			found = true
		}
	}

	if !found {
		if msg, err := images.Pull(Service.Config.Image); err != nil {
			return err
		} else if strings.Contains(msg, "already up to date") {
			return fmt.Errorf(msg)
		}
	} else {
		return fmt.Errorf("image already in registry, skipping")
	}

	return nil
}

// Start will perform a series of checks to see if the container starting
// is supposed be removed before-hand and will check to see if the
// container is running before it is actually started.
func (Service *Service) Start() error {

	name, err := Service.GetFieldString("name")
	discrete, _ := Service.GetFieldBool("discrete")
	interactive, _ := Service.GetFieldBool("interactive")
	output, _ := Service.GetFieldBool("output")
	purpose, _ := Service.GetFieldString("purpose")

	if err != nil {
		return nil
	}

	s := false

	if !Service.HostConfig.AutoRemove {
		var e error
		s, e = Service.Status()
		if e != nil {
			return e
		}
	}

	if s && !Service.HostConfig.AutoRemove && !discrete {
		color.Print(Green(fmt.Sprintf("Already running %s\n", name)))
		return nil
	}

	if purpose == "addkeys" {
		if e := containers.Kill(name); e != nil {
			fmt.Sprintln(e)
		}
		if e := containers.Remove(name); e != nil {
			fmt.Sprintln(e)
		}
		if e := Service.Create(); e != nil {
			fmt.Sprintln(e)
		}
	}

	if !interactive {
		err = Service.DockerRun()
		if err != nil {
			return err
		}

		l, _ := Service.DockerLogs()
		if output && string(l) != "" {
			fmt.Println(string(l))
		}

		if c, err := Service.ID(); c != "" {
			return nil
		} else if err != nil {
			return err
		}
	} else {
		err = Service.DockerRunInteractive()
		if err != nil {
			return err
		}
	}

	return nil
}

// Create will perform a series of checks to see if the container starting
// is supposed be removed before-hand and will check to see if the
// container is running before it is actually started.
func (Service *Service) Create() error {

	name, err := Service.GetFieldString("name")
	output, _ := Service.GetFieldBool("output")

	if err != nil || name == "" {
		return fmt.Errorf("missing name property")
	}

	err = Service.DockerCreate()
	if err != nil {
		return err
	}

	l, _ := Service.DockerLogs()
	if output && string(l) != "" {
		fmt.Println(string(l))
	}

	if c, err := Service.ID(); c != "" {
		return err
	}

	return nil
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
	containers, _ := containers.List()
	for _, container := range containers {
		for _, n := range container.Names {
			if strings.Contains(n, name) && strings.HasPrefix(container.Status, "Up") {
				return true, nil
			}
		}
	}

	return false, nil

}

// ID will get a types.Container variable for a given running container
// and it will not retrieve any information on containers that are not running.
// todo: digests instead of name?
func (Service *Service) ID() (string, error) {
	containers, _ := containers.List()
	for _, container := range containers {
		if _, ok := container.Labels["pygmy.name"]; ok {
			if strings.Contains(container.Labels["pygmy.name"], Service.Config.Labels["pygmy.name"]) {
				return container.ID, nil
			}
		}
	}
	return "", fmt.Errorf("container using image '%v' was not found\n", Service.Config.Image)
}

// Labels will get a types.Container variable for a given running container
// and it will not retrieve any information on containers that are not running.
func (Service *Service) Labels() (map[string]string, error) {
	containers, _ := containers.List()
	for _, container := range containers {
		if _, ok := container.Labels["pygmy.name"]; ok {
			if strings.Contains(container.Labels["pygmy.name"], Service.Config.Labels["pygmy.name"]) {
				return container.Labels, nil
			}
		}
	}
	return nil, fmt.Errorf("container using image '%v' was not found\n", Service.Config.Image)
}

// Clean will cleanup and remove the container.
func (Service *Service) Clean() error {

	pygmy, _ := Service.GetFieldBool("pygmy.enable")
	name, e := Service.GetFieldString("name")
	if e != nil {
		return nil
	}

	Containers, _ := containers.List()
	for _, container := range Containers {
		if container.Names[0] == name {
			if pygmy {
				name := strings.TrimLeft(container.Names[0], "/")
				if e := containers.Kill(container.ID); e == nil {
					if !Service.HostConfig.AutoRemove {
						color.Print(Green(fmt.Sprintf("Successfully killed %s\n", name)))
					}
				}
				if e := containers.Stop(container.ID); e == nil {
					if !Service.HostConfig.AutoRemove {
						color.Print(Green(fmt.Sprintf("Successfully stopped %s\n", name)))
					}
				}
				if e := containers.Remove(container.ID); e != nil {
					if !Service.HostConfig.AutoRemove {
						color.Print(Green(fmt.Sprintf("Successfully removed %s\n", name)))
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

	id, err := Service.ID()
	if err != nil {
		if !discrete {
			color.Print(Red(fmt.Sprintf("Not running %s\n", name)))
		}
		return nil
	}

	if e := containers.Stop(id); e == nil {
		if !discrete {
			containerName := strings.Trim(name, "/")
			color.Print(Green(fmt.Sprintf("Successfully stopped %v\n", containerName)))
		}
	}

	return nil
}

// StopAndRemove will stop and remove the container.
func (Service *Service) StopAndRemove() error {

	name, e := Service.GetFieldString("name")
	discrete, _ := Service.GetFieldBool("discrete")
	if e != nil {
		return nil
	}

	id, err := Service.ID()
	if err != nil {
		if !discrete {
			color.Print(Red(fmt.Sprintf("Not running %v\n", id)))
		}
		return nil
	}

	if e := containers.Stop(id); e == nil {
		if e := containers.Remove(id); e == nil {
			if !discrete {
				containerName := strings.Trim(name, "/")
				fmt.Print(Green(fmt.Sprintf("Successfully removed %v\n", containerName)))
			}
		}
	} else {
		return e
	}

	return nil
}

// Remove will stop the container.
func (Service *Service) Remove() error {

	discrete, _ := Service.GetFieldBool("discrete")
	id, _ := Service.ID()

	containerName := strings.Trim(id, "/")
	if e := containers.Remove(id); e == nil {
		if !discrete {
			fmt.Print(Green(fmt.Sprintf("Successfully removed %s\n", containerName)))
		}
	} else {
		return e
	}

	return nil
}

/// BELOW IS NOT IN SPEC TO THE INTERFACE.

// DockerLogs will return the logs from the container.
func (Service *Service) DockerLogs() ([]byte, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	cli.NegotiateAPIVersion(ctx)
	if err != nil {
		return []byte{}, err
	}

	name, _ := Service.GetFieldString("name")
	return containers.Logs(name)
}

// DockerRun will start an existing container.
func (Service *Service) DockerRun() error {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	cli.NegotiateAPIVersion(ctx)
	if err != nil {
		return err
	}

	name, e := Service.GetFieldString("name")
	if e != nil {
		return fmt.Errorf("container config is missing label for name")
	}
	if err := containers.Start(name, container.StartOptions{}); err != nil {
		return err
	}

	return nil

}

// DockerRunInteractive will start an interactive container.
func (Service *Service) DockerRunInteractive() error {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	cli.NegotiateAPIVersion(ctx)
	if err != nil {
		return err
	}

	name, e := Service.GetFieldString("name")
	if e != nil {
		return fmt.Errorf("container config is missing label for name")
	}

	waiter, err := containers.Attach(name, container.AttachOptions{
		Stderr: true,
		Stdout: true,
		Stdin:  true,
		Stream: true,
	})
	if err != nil {
		return err
	}

	// Connect the stdin/stdout/stderr streams to the container.
	go func() {
		if _, err := io.Copy(os.Stdout, waiter.Reader); err != nil {
			panic(fmt.Sprintf("Error streaming Stdout: %s", err))
		}
	}()

	go func() {
		if _, err := io.Copy(os.Stderr, waiter.Reader); err != nil {
			panic(fmt.Sprintf("Error streaming Stderr: %s", err))
		}
	}()

	go func() {
		if _, err := io.Copy(waiter.Conn, os.Stdin); err != nil {
			panic(fmt.Sprintf("Error streaming Stdin: %s", err))
		}
	}()

	if err := containers.Start(name, container.StartOptions{}); err != nil {
		return err
	}

	// Manipulate the terminal raw mode to support passing password prompts.
	fd := int(os.Stdin.Fd())
	var oldState *term.State
	if term.IsTerminal(fd) {
		oldState, err = term.MakeRaw(fd)
		if err != nil {
			return err
		}

		defer func() {
			if err := term.Restore(fd, oldState); err != nil {
				panic(fmt.Sprintf("Error restoring terminal: %s", err))
			}
		}()
	}

	if err := containers.Wait(name, container.WaitConditionNotRunning); err != nil {
		return err
	}

	return nil

}

// DockerCreate will setup and run a given container.
func (Service *Service) DockerCreate() error {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	cli.NegotiateAPIVersion(ctx)
	if err != nil {
		return err
	}

	// Sanity check to ensure we don't get name conflicts.
	c, _ := containers.List()
	for _, cn := range c {
		if strings.HasSuffix(cn.Names[0], Service.Config.Labels["pygmy.name"]) {
			return fmt.Errorf("container already created, or namespace is already taken")
		}
	}

	name, e := Service.GetFieldString("name")
	if e != nil {
		return fmt.Errorf("container config is missing label for name")
	}

	_, err = containers.Create(name, Service.Config, Service.HostConfig, Service.NetworkConfig)
	if err != nil {
		return err
	}

	return nil

}
