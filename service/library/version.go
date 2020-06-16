package library

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var (

	// isPygmy will determine if the local git branch is Pygmy.
	// otherwise the mechanisms to allow non-pygmy hash refs
	// being shown.
	isPygmy = false

	// Version tag used on local builds in Actions.
	COMMITSHA = ""

	// Fixed version which is modified via github actions for the
	// release builds. Should match version with v prepended.
	COMMITTAG = ""

	// Custom indicator will be added if changes are detected
	// to pygmy, so it would read dev-xxxxxxx-custom
	CUSTOMTAG = ""
)

// Version describes which version of Pygmy is running. This will be kept
// up to date as possible, and should be included in a release tag on the
// master branch, and should be changed to something adequate immediately
// after the release is published.
func Version(c Config) {

	// The Commit SHA is provided to GitHub Actions, if it is not set when compiled
	// we assume git is being used. Let's grab metadata from the repository for use
	// in the runtime -
	if COMMITSHA == "" {
		r, _ := exec.Command("git", "remote", "-v").Output()
		remotes := strings.Split(string(r), "\n")
		for remote := range remotes {
			if strings.Contains(remotes[remote], "pygmy-go") {
				isPygmy = true
			}
		}

		sha, _ := exec.Command("git", "rev-parse", "HEAD").Output()
		COMMITSHA = string(sha)

		_, changes := exec.Command("git", "diff-index", "--quiet", "HEAD").Output()
		if changes != nil {
			CUSTOMTAG = "-custom"
		}
	}

	// Detect a tagged version:
	if !isPygmy && COMMITTAG != "" {
		b := strings.Split(COMMITTAG, "/")
		reference := b[len(b)-1]

		if match, _ := regexp.Match("^v[0-9]+.[0-9]+.[0-9]+$", []byte(reference)); match {
			fmt.Printf("Pygmy %v%v\n", reference, CUSTOMTAG)
			return
		}
	}

	// Detect a SHA reference:
	if isPygmy || COMMITSHA != "" {
		if match, _ := regexp.Match("^[0-9|a-z|A-Z]+$", []byte(COMMITSHA)); match {
			fmt.Printf("Pygmy version dev-%v%v\n", COMMITSHA[0:7], CUSTOMTAG)
			return
		}
	} else {
		fmt.Printf("Pygmy version unidentifiable.\n")
	}
}
