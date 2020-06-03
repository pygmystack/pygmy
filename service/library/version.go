package library

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	// Version tag used on local builds in Travis.
	VERSIONTAG = os.Getenv("TRAVIS_COMMIT")
)

const (
	// Fixed version which is modified via travis for the
	// release builds. Should match version with v prepended.
	RELEASETAG = ""
)

// Version describes which version of Pygmy is running. This will be kept
// up to date as possible, and should be included in a release tag on the
// master branch, and should be changed to something adequate immediately
// after the release is published.
func Version(c Config) {

	// RELEASETAG will be provided via `sed` in the build pipeline.
	if RELEASETAG != "" {
		if match, _ := regexp.Match("^v[0-9]+.[0-9]+.[0-9]+$", []byte(RELEASETAG)); match {
			fmt.Printf("Pygmy %v\n", RELEASETAG)
			return
		} else if match, _ := regexp.Match("^[0-9|a-z|A-Z]+$", []byte(RELEASETAG)); match {
			fmt.Printf("Pygmy version dev-%v\n", RELEASETAG[0:7])
			return
		}
	}

	// Get tags and reference information.
	tags, _ := exec.Command("git", "show-ref", "--tags").Output()
	ref, _ := exec.Command("git", "rev-parse", "HEAD").Output()

	// Scan the references from the tags to check if the current reference is
	// associated to a tag. This is to show the tag version on a local build
	// from source.
	for _, tag := range strings.Split(string(tags), "\n") {
		if strings.Contains(tag, string(ref)) {
			fmt.Printf("Pygmy %v", strings.SplitAfter(tag, "/")[2])
			return
		}
	}

	if VERSIONTAG != "" {
		// If the version tag isn't empty:
		fmt.Printf("Pygmy %v", VERSIONTAG)
	} else {
		// If we don't have a version tag, use a reference.
		fmt.Printf("Pygmy version dev-%v\n", string(ref[0:7]))
	}
}
