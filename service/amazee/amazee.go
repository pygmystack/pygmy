package amazee

import (
	"github.com/docker/docker/api/types"
	model "github.com/fubarhouse/pygmy-go/service/interface"
)

// pull will perform an image update for a single image
// which is provided as a container provided by the
// Docker API.
func pull(image string) (string, error) {
	return model.DockerPull(image)
}

// list will return all running containers,
// equivalent to a `docker ps` command.
func list() ([]types.ImageSummary, error) {
	images, err := model.DockerImageList()
	return images, err
}
