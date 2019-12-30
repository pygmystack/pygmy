package model

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// DockerVolumeExists will check if a Docker volume exists.
func DockerVolumeExists(name string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return false, err
	}
	_, _, err = cli.VolumeInspectWithRaw(ctx, name)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DockerVolumeCreate will create a Docker volume.
func DockerVolumeCreate(name string) (types.Volume, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return types.Volume{}, err
	}
	return cli.VolumeCreate(ctx, volume.VolumesCreateBody{
		Name: name,
	})
}
