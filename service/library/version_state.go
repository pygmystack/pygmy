package library

import "fmt"

// PYGMY_VERSION is the equivalent to the version pygmy is being associated to.
// This variable is exclusively used when packaging a formal release.
var PYGMY_VERSION = ""

func printversion() bool {
	if PYGMY_VERSION == "" {
		return false
	}
	fmt.Printf("Pygmy version v%v\n", PYGMY_VERSION)
	return true
}
