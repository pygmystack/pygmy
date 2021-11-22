package library

import (
	"fmt"
	"strings"
)

// PYGMY_VERSION is the equivalent to the version pygmy is being associated to.
// This variable is exclusively used when packaging a formal release.
var PYGMY_VERSION = ""

func printversion() bool {
	parts := strings.Split(PYGMY_VERSION, "/")
	if PYGMY_VERSION == "" {
		return false
	}
	resultVersion :=  fmt.Sprintf("%v", parts[len(parts)-1])
	if !strings.HasPrefix(resultVersion, "v") {
		resultVersion = fmt.Sprintf("v%v", resultVersion)
	}
	fmt.Printf("Pygmy version %v", resultVersion)
	return true
}