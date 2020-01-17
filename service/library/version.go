package library

import "fmt"

// Version describes which version of Pygmy is running. This will be kept
// up to date as possible, and should be included in a release tag on the
// master branch, and should be changed to something adequate immediately
// after the release is published.
func Version(c Config) {
	fmt.Println("version called")
}
