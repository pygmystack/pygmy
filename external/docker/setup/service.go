package setup

import (
	"fmt"

	"github.com/imdario/mergo"
	dockerruntime "github.com/pygmystack/pygmy/internal/runtime/docker"
)

// mergeService will merge two Service objects.
func mergeService(destination dockerruntime.Service, src *dockerruntime.Service) (*dockerruntime.Service, error) {
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
	return &destination, nil
}

// GetService will return a service from the configuration.
// This merges the information down to return the object, so it cannot be implemented in the another package.
func GetService(s dockerruntime.Service, c dockerruntime.Service) dockerruntime.Service {
	Service, _ := mergeService(s, &c)
	return *Service
}
