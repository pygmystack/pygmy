package amazee

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// AmazeeImagePull is the entrypoint for this module.
// It will trigger the image pull after identifying all
// the images which match the criteria.
func AmazeeImagePull() {
	pull_all()
}

// pull will perform an image update for a single image
// which is provided as a container provided by the
// Docker API.
func pull(container types.Container) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	_, err = cli.ImagePull(ctx, "docker.io/"+container.Image, types.ImagePullOptions{})
	if err != nil {
		fmt.Println(err)
	}
}

// list will return all running containers,
// equivelant to a `docker ps` command.
func list() ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet:   true,
	})
	return containers, err
}

// pull_all is a loop which will trigger a `docker pull` command
// for all images matching the criteria - using the Docker API.
func pull_all() {
	list, _ := list()
	for _, container := range list {
		if strings.Contains(container.Image, "amazeeio") {
			pull(container)
		}
	}
}