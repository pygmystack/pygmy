// Package setup is a package which exposes the commands externally to the compiled binaries.
package setup

import (
	networktypes "github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	dockerruntime "github.com/pygmystack/pygmy/internal/runtime/docker"
	"github.com/pygmystack/pygmy/internal/utils/resolv"
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
	Services map[string]dockerruntime.Service `yaml:"services"`

	SortedServices []string

	// Networks is for network configuration
	Networks map[string]networktypes.Inspect `yaml:"networks"`

	// NoDefaults will prevent default configuration items.
	Defaults bool

	// JSONFormat indicates the `status` command should print to stdout in JSON format.
	JSONFormat bool

	// JSONStatus contains JSON status content.
	JSONStatus StatusJSON

	// ResolversDisabled will disable the creation of any resolv configurations.
	ResolversDisabled bool

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
	URLValidations   []StatusJSONURLValidation   `json:"url_validations"`
}

type StatusJSONURLValidation struct {
	Endpoint string `json:"endpoint"`
	Success  bool   `json:"success"`
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
