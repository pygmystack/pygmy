package docker

import (
	"context"
	"github.com/docker/docker/client"
	containercontext "github.com/pygmystack/pygmy/internal/runtime/docker/docker/context"
	"os"
)

func NewClient() (*client.Client, context.Context, error) {
	ctx := context.Background()
	clientOpts := []client.Opt{
		client.WithAPIVersionNegotiation(),
	}
	if os.Getenv("DOCKER_HOST") != "" {
		clientOpts = append(clientOpts, client.FromEnv)
	} else if currentDockerHost, err := containercontext.CurrentDockerHost(); err != nil {
		return nil, nil, err
	} else if currentDockerHost != "" {
		clientOpts = append(clientOpts, client.WithHost(currentDockerHost))
	}
	cli, err := client.NewClientWithOpts(clientOpts...)
	if err != nil {
		return nil, nil, err
	}
	return cli, ctx, nil
}
