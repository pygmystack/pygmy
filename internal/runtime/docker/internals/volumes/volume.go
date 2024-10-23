package volumes

import (
	"github.com/docker/docker/api/types/volume"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
)

// Exists will check if a Docker volume has been created.
func Exists(volume string) (bool, error) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return false, err
	}
	_, _, err = cli.VolumeInspectWithRaw(ctx, volume)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Get will return the full contents of a types.Volume from the API.
func Get(name string) (volume.Volume, error) {
	cli, ctx, err := internals.NewClient()

	if err != nil {
		return volume.Volume{
			Name: name,
		}, err
	}

	volumes, err := cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return volume.Volume{
			Name: name,
		}, err
	}

	for _, volume := range volumes.Volumes {
		if volume.Name == name {
			return *volume, nil
		}
	}

	return volume.Volume{
		Name: name,
	}, nil
}

// Create will create a Docker Volume as configured.
func Create(volumeInput volume.Volume) (volume.Volume, error) {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return volume.Volume{}, err
	}
	return cli.VolumeCreate(ctx, volume.CreateOptions{
		Driver:     volumeInput.Driver,
		DriverOpts: volumeInput.Options,
		Labels:     volumeInput.Labels,
		Name:       volumeInput.Name,
	})
}
