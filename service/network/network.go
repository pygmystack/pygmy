package network

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	model "github.com/pygmystack/pygmy/service/interface"
)

// New will generate the defaults for the Docker network.
// If configuration is provided this will not be used at all.
func New(c *model.Params) types.NetworkResource {
	return types.NetworkResource{
		Name: fmt.Sprintf("%s-network", c.Prefix),
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
			"pygmy.name": fmt.Sprintf("%s-network", c.Prefix),
		},
	}
}
