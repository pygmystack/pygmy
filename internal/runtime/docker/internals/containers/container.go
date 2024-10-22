package containers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
	"io"
	"runtime"
	"strings"

	"github.com/containerd/platforms"
	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// Stop will stop the container.
func Stop(name string) error {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return err
	}
	timeout := 10
	err = cli.ContainerStop(ctx, name, containertypes.StopOptions{Timeout: &timeout})
	if err != nil {
		return err
	}
	return nil
}

// Kill will kill the container.
func Kill(name string) error {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return err
	}
	err = cli.ContainerKill(ctx, name, "")
	if err != nil {
		return err
	}
	return nil
}

// Remove will remove the container.
// It will not remove the image.
func Remove(id string) error {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return err
	}
	err = cli.ContainerRemove(ctx, id, containertypes.RemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Inspect will return the full container object.
func Inspect(container string) (types.ContainerJSON, error) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return types.ContainerJSON{}, err
	}

	return cli.ContainerInspect(ctx, container)
}

// Exec will run a command in a Docker container and return the output.
func Exec(container string, command string) ([]byte, error) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return []byte{}, err
	}

	rst, err := cli.ContainerExecCreate(ctx, container, containertypes.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          strings.Split(command, " ")})

	if err != nil {
		return []byte{}, err
	}

	response, err := cli.ContainerExecAttach(context.Background(), rst.ID, containertypes.ExecAttachOptions{})

	if err != nil {
		return []byte{}, err
	}

	data, _ := io.ReadAll(response.Reader)
	defer response.Close()
	return data, nil

}

// List will return a slice of containers
func List() ([]types.Container, error) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		fmt.Println(err)
	}

	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{
		All: true,
	})
	if err != nil {
		return []types.Container{}, err
	}

	return containers, nil

}

// Create will create a container, but will not run it.
func Create(ID string, config containertypes.Config, hostconfig containertypes.HostConfig, networkconfig networktypes.NetworkingConfig) (containertypes.CreateResponse, error) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return containertypes.CreateResponse{}, err
	}
	platform := platforms.Normalize(v1.Platform{
		Architecture: runtime.GOARCH,
		OS:           "linux",
	})
	resp, err := cli.ContainerCreate(ctx, &config, &hostconfig, &networkconfig, &platform, ID)
	if err != nil {
		return containertypes.CreateResponse{}, err
	}
	return resp, err
}

// Attach will return an attached response to a container.
func Attach(ID string, options containertypes.AttachOptions) (types.HijackedResponse, error) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return types.HijackedResponse{}, err
	}
	resp, err := cli.ContainerAttach(ctx, ID, options)
	if err != nil {
		return types.HijackedResponse{}, err
	}
	return resp, err
}

// Start will run an existing container.
func Start(ID string, options containertypes.StartOptions) error {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return err
	}
	if err := cli.ContainerStart(ctx, ID, containertypes.StartOptions{}); err != nil {
		return err
	}
	return err
}

// Wait will wait for the specificied container condition.
func Wait(ID string, condition containertypes.WaitCondition) error {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return err
	}
	statusCh, errCh := cli.ContainerWait(ctx, ID, condition)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
		return nil
	}

	return nil
}

// Logs will synchronously (blocking, non-concurrently) print
// logs to stdout and stderr, useful for quick containers with a small amount
// of output which are expected to exit quickly.
func Logs(ID string) ([]byte, error) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return []byte{}, err
	}
	b, e := cli.ContainerLogs(ctx, ID, containertypes.LogsOptions{
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
