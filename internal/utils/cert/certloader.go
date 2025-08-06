package cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	. "github.com/logrusorgru/aurora"
	"github.com/mitchellh/go-homedir"
	"github.com/pygmystack/pygmy/internal/utils/color"
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

		// If we have a cert, we verify it.
		if err := verifyCertificate(defaultCertPath); err != nil {
			return "", fmt.Errorf("failed to verify default TLS certificate: %w", err)
		}

		return defaultCertPath, nil
	}

	if flagCertPath != "" {
		// Check if the provided certificate path exists.
		if _, err := os.Stat(flagCertPath); os.IsNotExist(err) {
			return "", fmt.Errorf("TLS certificate file %s does not exist", flagCertPath)
		}
		// If we have a cert, we verify it.
		if err := verifyCertificate(flagCertPath); err != nil {
			return "", fmt.Errorf("failed to verify TLS certificate at %s: %w", flagCertPath, err)
		}
		return flagCertPath, nil
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
			break
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

	color.Print(Green(fmt.Sprintf("Successfully verified certificate and private key pair at %s.\n", certPath)))
	return nil
}
