package amazee

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	model "github.com/fubarhouse/pygmy-go/service/interface"
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
func pull(image string) error {
	return model.DockerPull(image)
}

// list will return all running containers,
// equivelant to a `docker ps` command.
func list() ([]types.ImageSummary, error) {
	images, err := model.DockerImageList()
	return images, err
}

// pull_all is a loop which will trigger a `docker pull` command
// for all images matching the criteria - using the Docker API.
func pull_all() {
	list, _ := list()
	for _, image := range list {
		if strings.Contains(fmt.Sprint(image.RepoTags), "amazeeio/") || strings.Contains(fmt.Sprint(image.RepoTags), "mailhog/mailhog") || strings.Contains(fmt.Sprint(image.RepoTags), "andyshinn/dnsmasq") {
			for _, tag := range image.RepoTags {
				err := pull(tag)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
