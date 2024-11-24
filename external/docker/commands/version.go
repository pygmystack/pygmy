package commands

import (
	"fmt"
	"github.com/pygmystack/pygmy/external/docker/setup"
	"runtime/debug"
)

// Version describes which version of Pygmy is running.
func Version(c setup.Config) {
	info, _ := debug.ReadBuildInfo()

	if info.Main.Version == "(devel)" {
		fmt.Println("Development version")
		return
	}

	fmt.Printf("Pygmy %s\n", info.Main.Version)
}
