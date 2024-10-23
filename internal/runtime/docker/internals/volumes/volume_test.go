package volumes

import (
	"context"
	"fmt"
	volume2 "github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
)

// testSetup will prepare the client for each test.
func testSetup() (context.Context, *client.Client) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		panic(err)
	}

	return ctx, cli
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

// TestCreate will test the API call for volume creation.
func TestCreate(t *testing.T) {
	ctx, cli := testSetup()
	volumeName := fmt.Sprintf("testVolume-%s", randomString(10))
	volume := volume2.Volume{
		Name: volumeName,
	}

	// Create a volume.
	_, err := Create(ctx, cli, volume)
	assert.NoError(t, err)

	// Remove the volume.
	err = Remove(ctx, cli, volumeName)
	assert.NoError(t, err)
}

// TestRemove will test the API call for volume removal.
func TestRemove(t *testing.T) {
	ctx, cli := testSetup()
	volumeName := fmt.Sprintf("testVolume-%s", randomString(10))
	volume := volume2.Volume{
		Name: volumeName,
	}

	// Create a volume.
	_, err := Create(ctx, cli, volume)
	assert.NoError(t, err)

	// Remove the volume.
	err = Remove(ctx, cli, volumeName)
	assert.NoError(t, err)
}

// TestGet will test the API call for volume getting.
func TestGet(t *testing.T) {
	ctx, cli := testSetup()
	volumeName := fmt.Sprintf("testVolume-%s", randomString(10))
	volume := volume2.Volume{
		Name: volumeName,
	}

	// Create a volume.
	_, err := Create(ctx, cli, volume)
	assert.NoError(t, err)

	// Get a volume
	v, err := Get(ctx, cli, volumeName)
	assert.NoError(t, err)
	assert.Equal(t, volumeName, v.Name)

	// Remove the volume.
	err = Remove(ctx, cli, volumeName)
	assert.NoError(t, err)
}

// TestExists will test the API call for volume existing.
func TestExists(t *testing.T) {
	ctx, cli := testSetup()
	volumeName := fmt.Sprintf("testVolume-%s", randomString(10))
	volume := volume2.Volume{
		Name: volumeName,
	}

	// Create a volume.
	_, err := Create(ctx, cli, volume)
	assert.NoError(t, err)

	// Check if the volume exists.
	exists, err := Exists(ctx, cli, volumeName)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Remove the volume.
	err = Remove(ctx, cli, volumeName)
	assert.NoError(t, err)
}
