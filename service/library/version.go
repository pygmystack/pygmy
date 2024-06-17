package library

import (
	"fmt"
	"runtime/debug"
)

// Version describes which version of Pygmy is running.
func Version(c Config) {
	info, _ := debug.ReadBuildInfo()

	if info.Main.Version == "(devel)" {
		fmt.Println("Development version")
		return
	}

	fmt.Printf("Pygmy v%s\n", info.Main.Version)
}
