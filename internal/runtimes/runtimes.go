package runtimes

import (
	img "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
)

// TODO.

type Runtime interface {
	ContainerAttach()
	ContainerCreate()
	ContainerInspect()
	ContainerKill()
	ContainerList()
	ContainerLogs()
	ContainerRemove()
	ContainerStart()
	ContainerStop()
	ContainerWait()

	ImageList() ([]img.Summary, error)
	ImagePull(image string) (string, error)

	NetworkCreate()
	NetworkRemove()
	NetworkConnect()
	NetworkDisconnect()

	NetworkList()
	NetworkGet()
	NetworkStatus()
	NetworkConnected()
	NetworkInspect()

	VolumeCreate(volumeInput volume.Volume) (volume.Volume, error)
	VolumeRemove()
	VolumeGet(name string) (volume.Volume, error)
	VolumeExists(volume string) (bool, error)
}
