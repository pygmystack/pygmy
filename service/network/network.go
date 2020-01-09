package network

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	model "github.com/fubarhouse/pygmy-go/service/interface"
)

// New will generate the defaults for the Docker network.
// If configuration is provided this will not be used at all.
func New() model.Network {
	return model.Network{
		Name:       "amazeeio-network",
		Containers: []string{"amazeeio-haproxy"},
		Config: types.NetworkCreate{
			IPAM: &network.IPAM{
				Config: []network.IPAMConfig{
					{
						Subnet:  "10.99.99.0/24",
						Gateway: "10.99.99.1",
					},
				},
			},
		},
	}
}
