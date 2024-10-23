package docker

import (
	containertypes "github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"

	"github.com/pygmystack/pygmy/internal/runtime"
)

type Service struct {
	runtime.Service
	Config        containertypes.Config
	HostConfig    containertypes.HostConfig
	Image         string `yaml:"image"`
	NetworkConfig networktypes.NetworkingConfig
}

// Params is an arbitrary struct to pass around configuration from the top
// level to the lowest level - such as variable input to one of the
// containers.
type Params struct {
	runtime.Params
	// Domain is the target domain for Pygmy to use.
	Domain string
}
