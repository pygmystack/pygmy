package model

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// ContainerRuntimeInterface is a WIP...
type ContainerRuntimeInterface interface {
	// ContainerCreate is problematic...
	ContainerCreate(ID string, config container.Config, hostconfig container.HostConfig, networkconfig network.NetworkingConfig) (container.ContainerCreateCreatedBody, error)
	// ContainerExec is agnostic.
	ContainerExec(container string, command string) ([]byte, error)
	// ContainerKill is agnostic.
	ContainerKill(name string) error
	// ContainerList is problematic.
	ContainerList() ([]types.Container, error)
	// ContainerLogs is agnostic.
	ContainerLogs(ID string) ([]byte, error)
	// ContainerRemove is agnostic.
	ContainerRemove(id string) error
	// ContainerStart is problematic.
	ContainerStart(ID string, options types.ContainerStartOptions) error
	// ContainerStop is agnostic.
	ContainerStop(name string) error
	// ImageList is problematic.
	ImageList() ([]types.ImageSummary, error)
	// ImagePull is agnostic.
	ImagePull(image string) (string, error)
	// NetworkCreate is problematic.
	NetworkCreate(network *types.NetworkResource) error
	// NetworkConnect is agnostic.
	NetworkConnect(network string, containerName string) error
	// NetworkCreate is agnostic.
	NetworkConnected(network string, containerName string) (bool, error)
	// NetworkGet is problematic.
	NetworkGet(name string) (types.NetworkResource, error)
	// NetworkRemove is agnostic.
	NetworkRemove(network string) error
	// NetworkStatus is agnostic.
	NetworkStatus(network string) (bool, error)
	// VolumeCreate is problematic.
	VolumeCreate(volume types.Volume) (types.Volume, error)
	// VolumeExists is problematic.
	VolumeExists(volume types.Volume) (bool, error)
	// VolumeGet is problematic.
	VolumeGet(name string) (types.Volume, error)
}
