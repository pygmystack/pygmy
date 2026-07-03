package cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"

	aur "github.com/logrusorgru/aurora"
	"github.com/mitchellh/go-homedir"

	"github.com/pygmystack/pygmy/internal/utils/color"
)

var ErrNoDefaultCertError = fmt.Errorf("no default TLS certificate path provided")

// GetDefaultCertPaths returns the default path for the TLS certificate.
func GetDefaultCertPaths() []string {
	homedir, _ := homedir.Dir()
	// Default certificate paths, can be overridden by user.
	return []string{
		path.Join(homedir, ".pygmy", "server.pem"),
		path.Join(homedir, "pygmy", "server.pem"),
	}
}

// ResolveCertPath checks if the provided TLS certificate path exists.
func ResolveCertPath(flagCertPath string) (string, error) {

	// Fetch the default paths to scan.
	defaultCertPaths := GetDefaultCertPaths()

	// Search default paths if no flag is inputted.
	if flagCertPath == "" {
		for _, defaultPath := range defaultCertPaths {
			if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
				continue
			}
			return defaultPath, nil
		}
	}

	// Assume a value has been provided, check for the certificate at the provided path.
	// If a custom value is provided, add it to the scan.
	mergedCertificatePaths := append(defaultCertPaths, flagCertPath)
	for _, filePath := range mergedCertificatePaths {
		// Check if the provided certificate path exists.
		if _, err := os.Stat(flagCertPath); os.IsNotExist(err) {
			return "", fmt.Errorf("TLS certificate file %s does not exist", flagCertPath)
		}
		// If we have a cert, we verify it.
		if err := verifyCertificate(flagCertPath); err != nil {
			return "", fmt.Errorf("failed to verify TLS certificate at %s: %w", flagCertPath, err)
		}
		return filePath, nil
	}

	return "", nil
}

// verifyCertificate will take a .pem and verify it is appropriate for use with HAProxy.
func verifyCertificate(certPath string) error {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("failed to read certificate file: %w", err)
	}

	var certs []*x509.Certificate
	var privateKey interface{}

	for {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		switch block.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse certificate: %w", err)
			}
			certs = append(certs, cert)
		case "PRIVATE KEY":
			privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse private key: %w", err)
			}
		case "RSA PRIVATE KEY":
			privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse RSA private key: %w", err)
			}
		case "EC PRIVATE KEY":
			privateKey, err = x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse EC private key: %w", err)
			}
		default:
			fmt.Println("Unknown PEM block type:", block.Type)
		}
	}
	if len(certs) == 0 {
		return fmt.Errorf("no certificates found in the provided file")
	}
	if privateKey == nil {
		return fmt.Errorf("no private key found in the provided file")
	}

	color.Print(aur.Green(fmt.Sprintf("Successfully verified certificate and private key pair at %s.\n", certPath)))
	return nil
}
