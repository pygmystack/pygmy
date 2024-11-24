package setup

import (
	"fmt"

	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/imdario/mergo"
)

// mergeVolume will merge two Volume objects.
func mergeVolume(destination volumetypes.Volume, src *volumetypes.Volume) (*volumetypes.Volume, error) {
	if err := mergo.Merge(&destination, src, mergo.WithOverride); err != nil {
		fmt.Println(err)
		return src, err
	}
	return &destination, nil
}

// GetVolume will return a volume from the configuration.
// This merges the information down to return the object, so it cannot be implemented in the volumes package.
func GetVolume(s volumetypes.Volume, c volumetypes.Volume) volumetypes.Volume {
	Volume, _ := mergeVolume(s, &c)
	return *Volume
}
