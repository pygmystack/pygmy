package library

import "fmt"

var PYGMY_VERSION = ""

func printversion() bool {
	if PYGMY_VERSION == "" {
		return false
	}
	fmt.Printf("Pygmy version v%v\n", PYGMY_VERSION)
	return true
}
