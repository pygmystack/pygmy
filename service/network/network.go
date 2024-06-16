package network

import (
	networktypes "github.com/docker/docker/api/types/network"
)

// New will generate the defaults for the Docker network.
// If configuration is provided this will not be used at all.
func New() networktypes.Inspect {
	return networktypes.Inspect{
		Name: "amazeeio-network",
		IPAM: networktypes.IPAM{
			Driver:  "",
			Options: nil,
			Config: []networktypes.IPAMConfig{
				{
					Subnet:  "10.99.99.0/24",
					Gateway: "10.99.99.1",
				},
			},
		},
		Labels: map[string]string{
			"pygmy.name": "amazeeio-network",
		},
	}
}
