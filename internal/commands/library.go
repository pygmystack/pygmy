// Package library is a package which exposes the commands externally to the compiled binaries.
package commands

import (
	"fmt"
	networktypes "github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/imdario/mergo"
	"github.com/pygmystack/pygmy/internal/runtime/docker"
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
	Services map[string]docker.Service `yaml:"services"`

	SortedServices []string

	// Networks is for network configuration
	Networks map[string]networktypes.Inspect `yaml:"networks"`

	// NoDefaults will prevent default configuration items.
	Defaults bool

	// JSONFormat indicates the `status` command should print to stdout in JSON format.
	JSONFormat bool

	// JSONStatus contains JSON status content.
	JSONStatus StatusJSON

	// Resolvers is for all resolvers
	Resolvers []resolv.Resolv `yaml:"resolvers"`

	// Volumes will ensure names volumes are created
	Volumes map[string]volumetypes.Volume
}

type StatusJSON struct {
	PortAvailability []string                    `json:"port_availability"`
	Services         map[string]StatusJSONStatus `json:"service_status"`
	Networks         []string                    `json:"networks"`
	Resolvers        []string                    `json:"resolvers"`
	Volumes          []string                    `json:"volumes"`
	SSHMessages      []string                    `json:"ssh_messages"`
	URLValidations   []string                    `json:"url_validations"`
}

type StatusJSONStatus struct {
	Container string `json:"container"`
	ImageRef  string `json:"image"`
	State     bool   `json:"running"`
}

// Key is a struct with SSH key details.
type Key struct {
	Path string `yaml:"path"`
}

func mergeService(destination docker.Service, src *docker.Service) (*docker.Service, error) {
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
	return &destination, nil
}

func getService(s docker.Service, c docker.Service) docker.Service {
	Service, _ := mergeService(s, &c)
	return *Service
}

func mergeNetwork(destination networktypes.Inspect, src *networktypes.Inspect) (*networktypes.Inspect, error) {
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
	return &destination, nil
}

func getNetwork(s networktypes.Inspect, c networktypes.Inspect) networktypes.Inspect {
	Network, _ := mergeNetwork(s, &c)
	return *Network
}

func mergeVolume(destination volumetypes.Volume, src *volumetypes.Volume) (*volumetypes.Volume, error) {
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
	return &destination, nil
}

func getVolume(s volumetypes.Volume, c volumetypes.Volume) volumetypes.Volume {
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
