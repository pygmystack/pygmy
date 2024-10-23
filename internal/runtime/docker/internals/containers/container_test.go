package containers

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/images"
)

func testSetup() (context.Context, *client.Client) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		panic(err)
	}

	images.Pull(ctx, cli, "nginx")

	return ctx, cli
}

func sampleData() (container.Config, container.HostConfig, network.NetworkingConfig) {
	config := container.Config{
		Image: "nginx",
	}
	hostConfig := container.HostConfig{}
	networkConfig := network.NetworkingConfig{}

	return config, hostConfig, networkConfig
}

func TestStop(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := "nginx-1337"

	// Create a container to stop.
	_, err := Create(ctx, cli, id, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Start the container to stop.
	err = Start(ctx, cli, id, container.StartOptions{})
	assert.NoError(t, err)

	// Inspect the container.
	dataBefore, err := Inspect(ctx, cli, id)
	assert.NoError(t, err)

	// Assert the container is running
	assert.Equal(t, dataBefore.State.Status, "running")

	// Stop the container.
	err = Stop(ctx, cli, id)
	assert.NoError(t, err)

	// Inspect the container.
	dataAfter, err := Inspect(ctx, cli, id)
	assert.NoError(t, err)

	// Assert the container is stopped.
	assert.Equal(t, dataAfter.State.Status, "exited")
	assert.NotNil(t, dataAfter.State.FinishedAt)

	// Clean up tests.
	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

func TestKill(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := "nginx-1336"

	// Create a container to kill.
	_, err := Create(ctx, cli, id, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Start the container to kill.
	err = Start(ctx, cli, id, container.StartOptions{})
	assert.NoError(t, err)

	// Inspect the container.
	dataBefore, err := Inspect(ctx, cli, id)
	assert.NoError(t, err)

	// Assert the container is running
	assert.Equal(t, dataBefore.State.Status, "running")

	// Kill the container.
	err = Kill(ctx, cli, id)
	assert.NoError(t, err)

	// Inspect the container.
	dataAfter, err := Inspect(ctx, cli, id)
	assert.NoError(t, err)

	// Assert the container is killed.
	assert.Equal(t, dataAfter.State.ExitCode, 137)

	// Clean up tests.
	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

func TestRemove(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := "nginx-1335"

	// Create a container to remove.
	_, err := Create(ctx, cli, id, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Stop the container.
	err = Stop(ctx, cli, id)
	assert.NoError(t, err)

	// Remove the container.
	err = Remove(ctx, cli, id)
	assert.NoError(t, err)

	// Inspect the container, expecting an error.
	_, err = Inspect(ctx, cli, id)
	assert.Error(t, err)
}

func TestInspect(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := "nginx-1334"

	// Create a container to kill.
	_, err := Create(ctx, cli, id, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Start the container to kill.
	err = Start(ctx, cli, id, container.StartOptions{})
	assert.NoError(t, err)

	// Inspect the container.
	dataAfter, err := Inspect(ctx, cli, id)
	assert.NoError(t, err)

	// Assert data is not nil
	assert.Equal(t, dataAfter.State.Status, "running")

	// Stop the container.
	err = Stop(ctx, cli, id)
	assert.NoError(t, err)

	// Clean up tests.
	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

func TestExec(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := "nginx-1332"

	// Container command
	config.Cmd = []string{"sh", "-c", "echo hello"}

	// Create a container to kill.
	_, err := Create(ctx, cli, id, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Start the container to kill.
	err = Start(ctx, cli, id, container.StartOptions{})
	assert.NoError(t, err)

	// Get the container logs
	data, err := Logs(ctx, cli, id)
	assert.NoError(t, err)

	// Assert the output
	assert.Contains(t, string(data), "hello")

	// Stop the container.
	err = Stop(ctx, cli, id)
	assert.NoError(t, err)

	// Clean up tests.
	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

func TestList(t *testing.T) {}

func TestCreate(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := "nginx-1333"

	// Create a container to kill.
	_, err := Create(ctx, cli, id, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Inspect the container.
	dataAfter, err := Inspect(ctx, cli, id)
	assert.NoError(t, err)

	// Assert the container is created.
	assert.Equal(t, dataAfter.State.Status, "created")

	// Clean up tests.
	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

func TestAttach(t *testing.T) {}

func TestStart(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := "nginx-1336"

	// Create a container to kill.
	_, err := Create(ctx, cli, id, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Start the container to kill.
	err = Start(ctx, cli, id, container.StartOptions{})
	assert.NoError(t, err)

	// Inspect the container.
	dataBefore, err := Inspect(ctx, cli, id)
	assert.NoError(t, err)

	// Assert the container is running
	assert.Equal(t, dataBefore.State.Status, "running")

	// Kill the container.
	err = Kill(ctx, cli, id)
	assert.NoError(t, err)

	// Inspect the container.
	dataAfter, err := Inspect(ctx, cli, id)
	assert.NoError(t, err)

	// Assert the container is killed.
	assert.Equal(t, dataAfter.State.ExitCode, 137)

	// Clean up tests.
	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

func TestWait(t *testing.T) {}

func TestLogs(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := "nginx-1331"

	// Container command
	config.Cmd = []string{"sh", "-c", "echo hello world"}

	// Create a container to kill.
	_, err := Create(ctx, cli, id, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Start the container to kill.
	err = Start(ctx, cli, id, container.StartOptions{})
	assert.NoError(t, err)

	// Get the container logs
	data, err := Logs(ctx, cli, id)
	assert.NoError(t, err)

	// Assert the output
	assert.Contains(t, string(data), "hello world")

	// Stop the container.
	err = Stop(ctx, cli, id)
	assert.NoError(t, err)

	// Clean up tests.
	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}
