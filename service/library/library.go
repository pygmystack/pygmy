// Package library is a package which exposes the commands externally to the compiled binaries.
package library

import (
	"fmt"

	"github.com/docker/docker/api/types/volume"

	"github.com/docker/docker/api/types"
	"github.com/imdario/mergo"

	model "github.com/pygmystack/pygmy/service/interface"
	"github.com/pygmystack/pygmy/service/resolv"
)

// Config is a struct of configurable options which can
// be passed to package library to configure logic for
// continued abstraction.
type Config struct {
	// Keys are the paths to the Keys which should be added.
	Keys []Key `yaml:"keys"`

	// Domain is the default domain suffix to use.
	Domain string `yaml:"domain"`

	// Services is a []model.Service for an index of all Services.
	Services map[string]model.Service `yaml:"services"`

	SortedServices []string

	// Networks is for network configuration
	Networks map[string]types.NetworkResource `yaml:"networks"`

	// NoDefaults will prevent default configuration items.
	Defaults bool

	// Resolvers is for all resolvers
	Resolvers []resolv.Resolv `yaml:"resolvers"`

	// Volumes will ensure names volumes are created
	Volumes map[string]volume.Volume
}

// Key is a struct with SSH key details.
type Key struct {
	Path string `yaml:"path"`
}

func mergeService(destination model.Service, src *model.Service) (*model.Service, error) {
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
	return &destination, nil
}

func getService(s model.Service, c model.Service) model.Service {
	Service, _ := mergeService(s, &c)
	return *Service
}

func mergeNetwork(destination types.NetworkResource, src *types.NetworkResource) (*types.NetworkResource, error) {
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
	return &destination, nil
}

func getNetwork(s types.NetworkResource, c types.NetworkResource) types.NetworkResource {
	Network, _ := mergeNetwork(s, &c)
	return *Network
}

func mergeVolume(destination volume.Volume, src *volume.Volume) (*volume.Volume, error) {
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
	return &destination, nil
}

func getVolume(s volume.Volume, c volume.Volume) volume.Volume {
	Volume, _ := mergeVolume(s, &c)
	return *Volume
}

// unique will return a slice with duplicates
// removed. It performs a similar function to
// the linux program `uniq`
func unique(stringSlice []string) []string {
	m := make(map[string]bool)
	for _, item := range stringSlice {
		if _, ok := m[item]; !ok {
			m[item] = true
		}
	}

	var result []string
	for item := range m {
		result = append(result, item)
	}
	return result
}
