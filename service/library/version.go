package library

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Version describes which version of Pygmy is running.
func Version(c Config) error {

	version = getVersion()
	fmt.Printf("Application version: %s\n", version)

	// Use the version information as needed
	log.Printf("Running version %s (commit: %s, built on: %s)", version, commit, date)

	return nil
}

func getVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value[:8] // Short commit hash
			case "vcs.time":
				date = setting.Value
			case "vcs.modified":
				if setting.Value == "true" {
					commit += "-dirty"
				}
			}
		}
	}

	if version == "dev" {
		return fmt.Sprintf("%s-%s", version, commit)
	}

	return strings.TrimPrefix(version, "v")
}
