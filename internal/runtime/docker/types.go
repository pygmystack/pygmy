package docker

import (
	containertypes "github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
)

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
	// TLSCertPath is the TLS Certificate Path.
	TLSCertPath string
}
