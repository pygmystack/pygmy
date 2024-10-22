package images

import (
	"encoding/json"
	"fmt"
	"github.com/pygmystack/pygmy/internal/runtimes/docker"
	"io"
	"regexp"
	"strings"

	img "github.com/docker/docker/api/types/image"
	"github.com/pygmystack/pygmy/service/endpoint"
)

// List will return a slice of Docker images.
func List() ([]img.Summary, error) {
	cli, ctx, err := docker.NewClient()
	if err != nil {
		fmt.Println(err)
	}

	images, err := cli.ImageList(ctx, img.ListOptions{
		All: true,
	})
	if err != nil {
		return []img.Summary{}, err
	}

	return images, nil

}

// Pull will pull a Docker image into the daemon.
func Pull(image string) (string, error) {
	cli, ctx, err := docker.NewClient()
	if err != nil {
		fmt.Println(err)
	}

	{

		// To support image references from external sources to docker.io we need to check
		// and validate the image reference for all known cases of validity.

		if m, _ := regexp.MatchString("^(([a-zA-Z0-9_.-]+)[/]([a-zA-Z0-9_.-]+)[/]([a-zA-Z0-9_.-]+)[:]([a-zA-Z0-9_.-]+))$", image); m {
			// URL was provided (in full), but the tag was provided.
			// For this, we do not alter the value provided.
			// Examples:
			//  - quay.io/pygmystack/pygmy:latest
			image = fmt.Sprintf("%v", image)
		} else if m, _ := regexp.MatchString("^(([a-zA-Z0-9_.-]+)[/]([a-zA-Z0-9_.-]+)[/]([a-zA-Z0-9_.-]+))$", image); m {
			// URL was provided (in full), but the tag was not provided.
			// For this, we do not alter the value provided.
			// Examples:
			//  - quay.io/pygmystack/pygmy
			image = fmt.Sprintf("%v:latest", image)
		} else if m, _ := regexp.MatchString("^(([a-zA-Z0-9_.-]+)[/]([a-zA-Z0-9_.-]+)[:]([a-zA-Z0-9_.-]+))$", image); m {
			// URL was not provided (in full), but the tag was provided.
			// For this, we prepend 'docker.io/' to the reference.
			// Examples:
			//  - pygmystack/pygmy:latest
			image = fmt.Sprintf("docker.io/%v", image)
		} else if m, _ := regexp.MatchString("^(([a-zA-Z0-9_.-]+)[/]([a-zA-Z0-9_.-]+))$", image); m {
			// URL was not provided (in full), but the tag was not provided.
			// For this, we prepend 'docker.io/' to the reference.
			// Examples:
			//  - pygmystack/pygmy
			image = fmt.Sprintf("docker.io/%v:latest", image)
		} else if m, _ := regexp.MatchString("^(([a-zA-Z0-9_.-]+)[:]([a-zA-Z0-9_.-]+))$", image); m {
			// Library image was provided with tag identifier.
			// For this, we prepend 'docker.io/' to the reference.
			// Examples:
			//  - pygmy:latest
			image = fmt.Sprintf("docker.io/%v", image)
		} else if m, _ := regexp.MatchString("^([a-zA-Z0-9_.-]+)$", image); m {
			// Library image was provided without tag identifier.
			// For this, we prepend 'docker.io/' to the reference.
			// Examples:
			//  - pygmy
			image = fmt.Sprintf("docker.io/%v:latest", image)
		} else {
			// Validation not successful
			return "", fmt.Errorf("error: regexp validation for %v failed", image)
		}
	}

	// DockerHub Registry causes a stack trace fatal error when unavailable.
	// We can check for this and report back, handling it gracefully and
	// tell the user the service is down momentarily, and to try again shortly.
	if strings.HasPrefix(image, "docker.io") {
		if s := endpoint.Validate("https://registry-1.docker.io/v2/"); !s {
			return "", fmt.Errorf("cannot reach the Docker Hub Registry, please try again in a few minutes")
		}
	}

	data, err := cli.ImagePull(ctx, image, img.PullOptions{})
	d := json.NewDecoder(data)

	type Event struct {
		Status         string `json:"status"`
		Error          string `json:"error"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"progressDetail"`
	}

	var event *Event
	if err == nil {
		for {
			if err := d.Decode(&event); err != nil {
				if err == io.EOF {
					break
				}

				panic(err)
			}
		}

		if event != nil {
			if strings.Contains(event.Status, "Downloaded newer image") {
				return fmt.Sprintf("Successfully pulled %v", image), nil
			}

			if strings.Contains(event.Status, "Image is up to date") {
				return fmt.Sprintf("Image %v is up to date", image), nil
			}
		}

		return event.Status, nil
	}

	if strings.Contains(err.Error(), "pull access denied") {
		return fmt.Sprintf("Error trying to update image %v: pull access denied", image), nil
	}

	return "", nil
}
