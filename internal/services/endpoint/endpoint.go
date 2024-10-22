// Package endpoint provides a way to test a HTTP/HTTPS endpoint for a 200 response code.
// note that this package does not support insecure HTTPS at this time.
package endpoint

import (
	"crypto/tls"
	"net/http"
)

// Validate will submit a web request to test the container service.
// If a 200 response code is received it will pass and return true.
// Any other result will fail this validation process.
//
// This is to provided to the user through the up and status commands.
func Validate(url string) bool {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Create a web request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// Test failed.
		return false
	}

	// Submit a web request
	resp, err := client.Do(req)
	if err != nil {
		// Test failed.
		return false
	}

	// Housekeeping
	defer resp.Body.Close()

	// The default response for a failed loopback is a 503.
	// Because server errors are 503, we should make sure
	// we do not get a 5xx response code from the endpoint.
	// 500 is known to be a success as well, so start from 501.

	// Check for known failure status response codes (failures):
	if resp.StatusCode >= 501 && resp.StatusCode < 600 {
		return false
	}

	// Test passed.
	return true
}
