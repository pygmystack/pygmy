package runtimes

import (
	containertypes "github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/container"
)

// TODO.

type Service struct {
	Config        containertypes.Config
	HostConfig    containertypes.HostConfig
	Image         string `yaml:"image"`
	NetworkConfig networktypes.NetworkingConfig
}

// Params is an arbitrary struct to pass around configuration from the top
// level to the lowest level - such as variable input to one of the
// containers.
type Params struct {
	// Domain is the target domain for Pygmy to use.
	Domain string
}

// ServiceRuntime is ...
type ServiceRuntime interface {
	Setup() error
	Start() error
	Create() error
	Status() (bool, error)
	GetRunning() (container.Container, error) // TODO migrate elsewhere to remove typed response.
	Clean() error
	Stop() error
	StopAndRemove() error
	Remove() error
}
