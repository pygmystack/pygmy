package amazee

import (
	"fmt"
	"strings"

	model "github.com/fubarhouse/pygmy/service/interface"
)

func AmazeeImagePull() {
	pull_all()
}

func pull(image_name string) {

	Output, Error := model.DockerRun([]string{"pull", image_name})

	if Error != nil {
		model.Red(fmt.Sprintf("Failed to update %v. Command 'docker pull %v' failed", image_name, image_name))
		return
	}

	if strings.Contains(string(Output), "Image is up to date") {
		model.Green(fmt.Sprintf("Image %v is already up to date", image_name))
	} else {
		model.Green(fmt.Sprintf("Image %v was updated successfully", image_name))
	}

}

func ls_cmd() ([]string, error) {

	List, Error := model.DockerRun([]string{"image", "ls", "--format", "{{.Repository}}:{{.Tag}}"})
	if Error != nil {
		return []string{}, Error
	}
	// For better handling of containers, we should compare our
	// results against a whitelist instead of preferential
	// treatment of linux pipes.
	containers := strings.Split(string(List), "\n")
	amazeeContainers := []string{}
	for _, container := range containers {
		// Selectively target amazeeio/* images.
		if strings.Contains(container, "amazeeio/") {
			// Filter out items which we don't want.
			//if !strings.Contains(container, "none") || strings.Contains(container, "ssh-agent") || strings.Contains(container, "haproxy") {
				amazeeContainers = append(amazeeContainers, container)
			//}
		}
	}
	return amazeeContainers, nil

}

func pull_all() {

	list, _ := ls_cmd()
	for _, image := range list {
		pull(image)
	}

}