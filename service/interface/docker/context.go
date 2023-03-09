package docker

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type DockerConfig struct {
	CurrentContext string `json:"currentContext"`
}

type DockerContextManifest struct {
	Name      string
	Endpoints map[string]struct {
		Host string
	}
}

func filePathInHomeDir(elem ...string) (string, error) {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(append([]string{home}, elem...)...), nil
}

func CurrentContext() (string, error) {
	configPath, err := filePathInHomeDir(".docker", "config.json")
	if err != nil {
		return "", err
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", errors.New(err.(*os.PathError).Error())
	}

	dockerConfig := DockerConfig{}
	err = json.Unmarshal(configBytes, &dockerConfig)
	if err != nil {
		return "", err
	}

	return dockerConfig.CurrentContext, nil
}

func EndpointFromContext(context string) (string, error) {
	manifestDir, err := filePathInHomeDir(".docker", "contexts", "meta")
	if err != nil {
		return "", err
	}

	var contextManifest DockerContextManifest

	err = filepath.WalkDir(manifestDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() != "meta.json" {
			return nil
		}
		manifestBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		manifest := DockerContextManifest{}
		err = json.Unmarshal(manifestBytes, &manifest)
		if err != nil {
			return err
		}

		if manifest.Name == context {
			contextManifest = manifest
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return contextManifest.Endpoints["docker"].Host, nil
}

func CurrentDockerHost() (string, error) {
	dockerCurrentContext, err := CurrentContext()
	if err != nil {
		return "", err
	}

	currentDockerHost := ""
	if dockerCurrentContext != "" {
		currentDockerHost, err = EndpointFromContext(dockerCurrentContext)
		if err != nil {
			return "", err
		}
	}
	return currentDockerHost, nil
}
