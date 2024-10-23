package commands

import (
	"fmt"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
)

// Down will bring pygmy down safely
func Down(c Config) error {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return err
	}

	Setup(ctx, cli, &c)
	for _, Service := range c.Services {
		enabled, _ := Service.GetFieldBool(ctx, cli, "enable")
		if enabled {
			e := Service.StopAndRemove(ctx, cli)
			if e != nil {
				name, _ := Service.GetFieldString(ctx, cli, "name")
				fmt.Printf("Failed to stop and remove %s\n", name)
			}
		}
	}

	return nil
}
