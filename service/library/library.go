// Library is a package which exposes the commands externally to the compiled binaries.
package library

import (
	"fmt"

	"github.com/docker/docker/api/types"
	model "github.com/fubarhouse/pygmy-go/service/interface"
	"github.com/fubarhouse/pygmy-go/service/resolv"
	"github.com/imdario/mergo"
)

// Config is a struct of configurable options which can
// be passed to package library to configure logic for
// continued abstraction.
type Config struct {
	// Keys are the paths to the Keys which should be added.
	Keys []string `yaml:"Keys"`

	// SkipKey indicates key adding should be skipped.
	SkipKey bool

	// Services is a []model.Service for an index of all Services.
	Services map[string]model.Service `yaml:"services"`

	SortedServices []string

	// Networks is for network configuration
	Networks map[string]types.NetworkResource`yaml:"networks"`

	// NoDefaults will prevent default configuration items.
	Defaults bool

	// Resolvers is for all resolvers
	Resolvers []resolv.Resolv `yaml:"resolvers"`

	// Volumes will ensure names volumes are created
	Volumes []string
}

func mergeService(destination model.Service, src *model.Service) (*model.Service, error) {
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
	return &destination, nil
}

func getService(s model.Service, c model.Service) model.Service {
	Service, _ := mergeService(s, &c)
	return *Service
}
