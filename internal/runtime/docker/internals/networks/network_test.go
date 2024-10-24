package networks

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/containers"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/images"
)

// sampleData will generate sample data for use in each test.
func sampleData() (container.Config, container.HostConfig, network.NetworkingConfig) {
	config := container.Config{
		Image: "nginx",
	}
	hostConfig := container.HostConfig{}
	networkConfig := network.NetworkingConfig{}

	return config, hostConfig, networkConfig
}

// randomString will generate a random string.
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

// testSetup will prepare the client for each test.
func testSetup() (context.Context, *client.Client) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		panic(err)
	}

	return ctx, cli
}

// TestCreate will test the network creation API request.
func TestCreate(t *testing.T) {
	ctx, cli := testSetup()
	id := fmt.Sprintf("testNetwork-%s", randomString(10))

	network := &network.Inspect{
		Name: id,
		Labels: map[string]string{
			"pygmy.name": id,
		},
	}

	err := Create(ctx, cli, network)
	assert.NoError(t, err)

	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

// TestRemove will test the network removal API request.
func TestRemove(t *testing.T) {
	ctx, cli := testSetup()
	id := fmt.Sprintf("testNetwork-%s", randomString(10))

	network := &network.Inspect{
		Name: id,
		Labels: map[string]string{
			"pygmy.name": id,
		},
	}

	err := Create(ctx, cli, network)
	assert.NoError(t, err)

	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

// TestStatus will test the network status API request.
func TestStatus(t *testing.T) {
	ctx, cli := testSetup()
	id := fmt.Sprintf("testNetwork-%s", randomString(10))

	network := &network.Inspect{
		Name: id,
		Labels: map[string]string{
			"pygmy.name": id,
		},
	}

	err := Create(ctx, cli, network)
	assert.NoError(t, err)

	status, err := Status(ctx, cli, id)
	assert.NoError(t, err)
	assert.True(t, status)

	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

// TestGet will test the network getter API request.
func TestGet(t *testing.T) {
	ctx, cli := testSetup()
	id := fmt.Sprintf("testNetwork-%s", randomString(10))

	network := &network.Inspect{
		Name: id,
		Labels: map[string]string{
			"pygmy.name": id,
		},
	}

	err := Create(ctx, cli, network)
	assert.NoError(t, err)

	obj, err := Get(ctx, cli, id)
	assert.NoError(t, err)
	assert.Equal(t, network.Name, obj.Name)

	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

// TestConnect will test the network connection API request.
func TestConnect(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()
	id := fmt.Sprintf("testNetwork-%s", randomString(10))
	containerName := fmt.Sprintf("testContainer-%s", randomString(10))

	config.Labels = map[string]string{
		"pygmy.name": containerName,
	}

	// Pull the image
	_, err := images.Pull(ctx, cli, config.Image)
	assert.NoError(t, err)

	// Create a container to stop.
	_, err = containers.Create(ctx, cli, containerName, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Start the container to stop.
	err = containers.Start(ctx, cli, containerName, container.StartOptions{})
	assert.NoError(t, err)

	network := &network.Inspect{
		Name: id,
		Labels: map[string]string{
			"pygmy.name": id,
		},
	}

	err = Create(ctx, cli, network)
	assert.NoError(t, err)

	err = Connect(ctx, cli, id, containerName)
	assert.NoError(t, err)

	// Start the container to stop.
	err = containers.Stop(ctx, cli, containerName)
	assert.NoError(t, err)

	// Start the container to stop.
	err = containers.Remove(ctx, cli, containerName)
	assert.NoError(t, err)

	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}

// TestConnect will test the validity of a network connection with a container.
func TestConnected(t *testing.T) {
	ctx, cli := testSetup()
	config, hostconfig, networkconfig := sampleData()
	id := fmt.Sprintf("testNetwork-%s", randomString(10))
	containerName := fmt.Sprintf("testContainer-%s", randomString(10))

	config.Labels = map[string]string{
		"pygmy.name": containerName,
	}

	// Pull the image
	_, err := images.Pull(ctx, cli, config.Image)
	assert.NoError(t, err)

	// Create a container to stop.
	_, err = containers.Create(ctx, cli, containerName, config, hostconfig, networkconfig)
	assert.NoError(t, err)

	// Start the container to stop.
	err = containers.Start(ctx, cli, containerName, container.StartOptions{})
	assert.NoError(t, err)

	network := &network.Inspect{
		Name: id,
		Labels: map[string]string{
			"pygmy.name": id,
		},
	}

	err = Create(ctx, cli, network)
	assert.NoError(t, err)

	err = Connect(ctx, cli, id, containerName)
	assert.NoError(t, err)

	connection, err := Connected(ctx, cli, id, containerName)
	assert.NoError(t, err)
	assert.True(t, connection)

	// Start the container to stop.
	err = containers.Stop(ctx, cli, containerName)
	assert.NoError(t, err)

	// Start the container to stop.
	err = containers.Remove(ctx, cli, containerName)
	assert.NoError(t, err)

	err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}
