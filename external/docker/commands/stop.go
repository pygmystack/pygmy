package commands

import (
	"fmt"

	"github.com/pygmystack/pygmy/external/docker/setup"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
)

// Stop will bring pygmy down safely
func Stop(c setup.Config) error {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return err
	}

	setup.Setup(ctx, cli, &c)

	for _, Service := range c.Services {
		enabled, _ := Service.GetFieldBool(ctx, cli, "enable")
		if enabled {
			e := Service.Stop(ctx, cli)
			if e != nil {
				fmt.Println(e)
			}
		}
	}

	return nil
}
