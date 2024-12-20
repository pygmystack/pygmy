package images

import (
	"context"
	"slices"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
)

func testSetup() (context.Context, *client.Client) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		panic(err)
	}

	return ctx, cli
}

// TestPullAndList will test the image functionality for pulling and listing a Docker image.
func TestPullAndList(t *testing.T) {
	ctx, cli := testSetup()

	id := "nginx:latest"

	// Remove the image from the registry.
	// We specifically do not want check this error.
	_, _ = Remove(ctx, cli, id)

	// Pull the image into the registry.
	pullResponse, err := Pull(ctx, cli, id)
	assert.NoError(t, err)

	// Ensure the output from this test contains some expected text.
	assert.Contains(t, pullResponse, "docker.io/nginx:latest")

	// List the images in the registry.
	list, err := List(ctx, cli)
	assert.NoError(t, err)

	// Check for the image in the registry.
	foundNginxImage := false
	for _, img := range list {
		if slices.Contains(img.RepoTags, "nginx:latest") {
			foundNginxImage = true
		}
	}
	assert.True(t, foundNginxImage)

	// Clean-up for this test.
	_, err = Remove(ctx, cli, id)
	assert.NoError(t, err)
}
