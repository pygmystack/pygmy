package cert

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
)

var ErrNoDefaultCertError = fmt.Errorf("no default TLS certificate path provided")

// GetDefaultCertPath returns the default path for the TLS certificate.
func GetDefaultCertPath() string {
	homedir, _ := homedir.Dir()
	// Default certificate path, can be overridden by user.
	return fmt.Sprintf("%v%vpygmy%vserver.pem", homedir, string(os.PathSeparator), string(os.PathSeparator))
}

// ResolveCertPath checks if the provided TLS certificate path exists.
func ResolveCertPath(flagCertPath string) (string, error) {

	defaultCertPath := GetDefaultCertPath()

	if flagCertPath == defaultCertPath { // If a default cert path is provided, check if it exists.
		// Check if the default certificate path exists.
		if _, err := os.Stat(defaultCertPath); os.IsNotExist(err) {
			return "", ErrNoDefaultCertError
		}

		return defaultCertPath, nil
	}

	if flagCertPath != "" {
		// Check if the provided certificate path exists.
		if _, err := os.Stat(flagCertPath); os.IsNotExist(err) {
			return "", fmt.Errorf("TLS certificate file %s does not exist", flagCertPath)
		}
		return flagCertPath, nil
	}

	return "", nil
}
