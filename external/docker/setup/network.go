package setup

import (
	"fmt"

	"github.com/imdario/mergo"

	networktypes "github.com/docker/docker/api/types/network"
)

// mergeNetwork will merge two Network objects.
func mergeNetwork(destination networktypes.Inspect, src *networktypes.Inspect) (*networktypes.Inspect, error) {
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
	return &destination, nil
}

// GetNetwork will return a network from the configuration.
// This merges the information down to return the object, so it cannot be implemented in the networks package.
func GetNetwork(s networktypes.Inspect, c networktypes.Inspect) networktypes.Inspect {
	Network, _ := mergeNetwork(s, &c)
	return *Network
}
