package containers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/client"
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
func Stop(ctx context.Context, client *client.Client, name string) error {
	timeout := 10
	err := client.ContainerStop(ctx, name, containertypes.StopOptions{Timeout: &timeout})
	if err != nil {
		return err
	}
	return nil
}

// Kill will kill the container.
func Kill(ctx context.Context, client *client.Client, name string) error {
	err := client.ContainerKill(ctx, name, "")
	if err != nil {
		return err
	}
	return nil
}

// Remove will remove the container.
// It will not remove the image.
func Remove(ctx context.Context, client *client.Client, id string) error {
	err := client.ContainerRemove(ctx, id, containertypes.RemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Inspect will return the full container object.
func Inspect(ctx context.Context, client *client.Client, container string) (types.ContainerJSON, error) {
	return client.ContainerInspect(ctx, container)
}

// Exec will run a command in a Docker container and return the output.
func Exec(ctx context.Context, client *client.Client, container string, command string) ([]byte, error) {
	rst, err := client.ContainerExecCreate(ctx, container, containertypes.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          strings.Split(command, " ")})

	if err != nil {
		return []byte{}, err
	}

	response, err := client.ContainerExecAttach(context.Background(), rst.ID, containertypes.ExecAttachOptions{})

	if err != nil {
		return []byte{}, err
	}

	data, _ := io.ReadAll(response.Reader)
	defer response.Close()
	return data, nil

}

// List will return a slice of containers
func List(ctx context.Context, client *client.Client) ([]types.Container, error) {
	containers, err := client.ContainerList(ctx, containertypes.ListOptions{
		All: true,
	})
	if err != nil {
		return []types.Container{}, err
	}

	return containers, nil
}

// Create will create a container, but will not run it.
func Create(ctx context.Context, client *client.Client, ID string, config containertypes.Config, hostconfig containertypes.HostConfig, networkconfig networktypes.NetworkingConfig) (containertypes.CreateResponse, error) {
	platform := platforms.Normalize(v1.Platform{
		Architecture: runtime.GOARCH,
		OS:           "linux",
	})
	resp, err := client.ContainerCreate(ctx, &config, &hostconfig, &networkconfig, &platform, ID)
	if err != nil {
		return containertypes.CreateResponse{}, err
	}
	return resp, err
}

// Attach will return an attached response to a container.
func Attach(ctx context.Context, client *client.Client, ID string, options containertypes.AttachOptions) (types.HijackedResponse, error) {
	resp, err := client.ContainerAttach(ctx, ID, options)
	if err != nil {
		return types.HijackedResponse{}, err
	}
	return resp, err
}

// Start will run an existing container.
func Start(ctx context.Context, client *client.Client, ID string, options containertypes.StartOptions) error {
	return client.ContainerStart(ctx, ID, containertypes.StartOptions{})
}

// Wait will wait for the specificied container condition.
func Wait(ctx context.Context, client *client.Client, ID string, condition containertypes.WaitCondition) error {
	statusCh, errCh := client.ContainerWait(ctx, ID, condition)
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
func Logs(ctx context.Context, client *client.Client, ID string) ([]byte, error) {
	b, e := client.ContainerLogs(ctx, ID, containertypes.LogsOptions{
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
