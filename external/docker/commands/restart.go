package commands

import (
	"fmt"
	"github.com/pygmystack/pygmy/external/docker/setup"
)

// Restart will stop and start Pygmy in its entirety.
func Restart(c setup.Config) {
	err := Down(c)
	if err != nil {
		fmt.Println(err)
	}

	err = Up(c)
	if err != nil {
		fmt.Println(err)
	}
}
