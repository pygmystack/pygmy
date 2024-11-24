package volumes

import (
	"context"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// Exists will check if a Docker volume has been created.
func Exists(ctx context.Context, cli *client.Client, volume string) (bool, error) {
	_, _, err := cli.VolumeInspectWithRaw(ctx, volume)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Get will return the full contents of a types.Volume from the API.
func Get(ctx context.Context, cli *client.Client, name string) (volume.Volume, error) {
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
func Create(ctx context.Context, cli *client.Client, volumeInput volume.Volume) (volume.Volume, error) {
	return cli.VolumeCreate(ctx, volume.CreateOptions{
		Driver:     volumeInput.Driver,
		DriverOpts: volumeInput.Options,
		Labels:     volumeInput.Labels,
		Name:       volumeInput.Name,
	})
}

// Remove will remove a Docker volume, which will be used exclusively for testing.
func Remove(ctx context.Context, cli *client.Client, volume string) error {
	return cli.VolumeRemove(ctx, volume, false)
}
