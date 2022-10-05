package model

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// Params is an arbitrary struct to pass around configuration from the top
// level to the lowest level - such as variable input to one of the
// containers.
type Params struct {
	// Domain is the target domain for Pygmy to use.
	Domain string
}

// DockerService is the requirements for a Docker container to be compatible.
// The Service struct is used to implement this interface, and individual
// variables of type Service can/have overwritten them when logic deems
// it necessary.
type DockerService interface {
	Setup() error
	Status() (bool, error)
	Start() error
	Stop() error
}

// Service is a collection of requirements for starting a container and
// provides a way for config of any container to be overridden and start
// fully compatible with Docker's API.
type Service struct {
	Config        container.Config
	HostConfig    container.HostConfig
	NetworkConfig network.NetworkingConfig
}

// Network is a struct containing the configuration of a single Docker network
// including some extra fields so that Pygmy knows how to interact with the
// desired outcome.
type Network struct {
	// Name is the name of the network, it is independent of the map key which
	// will be used to configure pygmy but this field should match the map key.
	Name string `yaml:"name"`
	// Containers is a []string which indicates the names of the containers
	// that need to be connected to this network.
	Containers []string `yaml:"containers"`
	// Config is the actual Network configuration for the Docker Network.
	// It is the Network creation configuration as provided by the Docker API.
	Config types.NetworkCreate `yaml:"config"`
}
