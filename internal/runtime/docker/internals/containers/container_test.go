package containers

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/images"
)

// testSetup will prepare the client for each test.
func testSetup() (context.Context, *client.Client) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		panic(err)
	}

	_, err = images.Pull(ctx, cli, "nginx")
	if err != nil {
		panic(err)
	}

	return ctx, cli
}

// randomString will generate a random string.
func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// sampleData will generate sample data for use in each test.
func sampleData() (container.Config, container.HostConfig, network.NetworkingConfig) {
	config := container.Config{
		Image: "nginx",
	}
	hostConfig := container.HostConfig{}
	networkConfig := network.NetworkingConfig{}

	return config, hostConfig, networkConfig
}

// TestStop will test the Stop operation of a container.
func TestStop(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := fmt.Sprintf("testContainer-%s", randomString(10))

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

// TestKill will test the Kill operation of a container.
func TestKill(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := fmt.Sprintf("testContainer-%s", randomString(10))

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

// TestRemove will test the Remove operation of a container.
func TestRemove(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := fmt.Sprintf("testContainer-%s", randomString(10))

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

// TestInspect will test the Inspect operation of a container.
func TestInspect(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := fmt.Sprintf("testContainer-%s", randomString(10))

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

// TestExec will test the Exec operation of a container.
func TestExec(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := fmt.Sprintf("testContainer-%s", randomString(10))

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

// TestList will test the List operation of a container.
func TestList(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := fmt.Sprintf("testContainer-%s", randomString(10))

	// Create a container to kill.
	_, err := Create(ctx, cli, id, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Start a container to kill.
	err = Start(ctx, cli, id, container.StartOptions{})
	assert.NoError(t, err)

	// Wait for the container.
	time.Sleep(time.Second * 3)

	// List the containers.
	list, err := List(ctx, cli)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	// Find the container
	foundContainer := false
	for _, ctr := range list {
		if strings.TrimPrefix(ctr.Names[0], "/") == id {
			foundContainer = true
		}
	}
	assert.True(t, foundContainer)

	// Stop the container.
	err = Stop(ctx, cli, id)
	assert.NoError(t, err)

	// Clean up tests.
	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

// TestCreate will test the Create operation of a container.
func TestCreate(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := fmt.Sprintf("testContainer-%s", randomString(10))

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

// TestAttach will test the Attach operation of a container.
func TestAttach(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := fmt.Sprintf("testContainer-%s", randomString(10))

	// Container command
	config.Cmd = []string{"sh", "-c", "echo 42"}

	// Create a container to kill.
	_, err := Create(ctx, cli, id, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Create a container to kill.
	err = Start(ctx, cli, id, container.StartOptions{})
	assert.NoError(t, err)

	attacher, err := Attach(ctx, cli, id, container.AttachOptions{
		Stream: true,
		Stdin:  false,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})

	// Read the output from the container.
	stdout, err := attacher.Reader.ReadString('\n')
	assert.NoError(t, err)

	// Check the output contains something we are expecting.
	assert.Contains(t, stdout, "42")

	// Create a container to kill.
	err = Stop(ctx, cli, id)
	assert.NoError(t, err)

	// Create a container to kill.
	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

// TestStart will test the Start operation of a container.
func TestStart(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := fmt.Sprintf("testContainer-%s", randomString(10))

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

// TestWait will test the Wait operation of a container.
func TestWait(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := fmt.Sprintf("testContainer-%s", randomString(10))

	// Container command
	config.Cmd = []string{"sh", "-c", "sleep 2 && exit 1"}

	// Create a container to kill.
	_, err := Create(ctx, cli, id, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Start the container to kill.
	err = Start(ctx, cli, id, container.StartOptions{})
	assert.NoError(t, err)

	err = Wait(ctx, cli, id, container.WaitConditionNotRunning)
	assert.NoError(t, err)

	// Stop the container.
	err = Stop(ctx, cli, id)
	assert.NoError(t, err)

	// Clean up tests.
	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

// TestLogs will test the Logs operation of a container.
func TestLogs(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()

	// Container ID.
	id := fmt.Sprintf("testContainer-%s", randomString(10))

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
