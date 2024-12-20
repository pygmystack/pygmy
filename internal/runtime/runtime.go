package runtime

import (
	"context"

	"github.com/docker/docker/client"
)

// ServiceRuntime is the definition of a Container Runtime for compatability with Pygmy.
type ServiceRuntime interface {
	Setup(ctx context.Context, cli *client.Client) error
	Start(ctx context.Context, cli *client.Client) error
	Create(ctx context.Context, cli *client.Client) error
	Status(ctx context.Context, cli *client.Client) (bool, error)
	Labels(ctx context.Context, cli *client.Client) (map[string]string, error)
	// @TODO: Does ID() work better as retrieving digests?
	ID(ctx context.Context, cli *client.Client) (string, error)
	Clean(ctx context.Context, cli *client.Client) error
	Stop(ctx context.Context, cli *client.Client) error
	StopAndRemove(ctx context.Context, cli *client.Client) error
	Remove(ctx context.Context, cli *client.Client) error

	SetField(ctx context.Context, cli *client.Client, name string, value interface{}) error
	GetFieldString(ctx context.Context, cli *client.Client, field string) (string, error)
	GetFieldInt(ctx context.Context, cli *client.Client, field string) (int, error)
	GetFieldBool(ctx context.Context, cli *client.Client, field string) (bool, error)
}
