// Package endpoint provides a way to test a HTTP/HTTPS endpoint for a 200 response code.
// note that this package does not support insecure HTTPS at this time.
package endpoint

import (
	"net/http"
)

// Validate will submit a web request to test the container service.
// If a 200 response code is received it will pass and return true.
// Any other result will fail this validation process.
//
// This is to provided to the user through the up and status commands.
func Validate(url string) bool {

	// Create a web request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// Test failed.
		return false
	}

	// Submit a web request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Test failed.
		return false
	}

	// Housekeeping
	defer resp.Body.Close()

	// Check for the desired result
	if resp.StatusCode == 200 {
		// Test passed!
		return true
	}

	// Test failed.
	return false
}
