package network

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

// New will generate the defaults for the Docker network.
// If configuration is provided this will not be used at all.
func New() types.NetworkResource {
	return types.NetworkResource{
		Name: "amazeeio-network",
		IPAM: network.IPAM{
			Driver:  "",
			Options: nil,
			Config: []network.IPAMConfig{
				{
					Subnet:  "10.99.99.0/24",
					Gateway: "10.99.99.1",
				},
			},
		},
		Labels: map[string]string{
			"pygmy": "pygmy",
		},
	}
}
