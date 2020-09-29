package library

import (
	"fmt"
)

// Version describes which version of Pygmy is running.
func Version(c Config) {

	// printversion is updated as static content via GitHub Actions.
	// If this version is not injected as static content, the version
	// is deemed unidentifiable - it should be assumed the binary was
	// compiled outside of official release management.
	if printversion() {
		return
	}

	fmt.Printf("Pygmy version unidentifiable.\n")
}
