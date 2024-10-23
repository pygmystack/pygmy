package context

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCurrentDockerHost will test the CurrentDockerHost function.
func TestCurrentDockerHost(t *testing.T) {
	_, err := CurrentDockerHost()
	assert.Nil(t, err)
}

// TestFilePathInHomeDir will test the filePathInHomeDir function.
func TestFilePathInHomeDir(t *testing.T) {
	filepath, err := filePathInHomeDir(".bash_profile")
	assert.NoError(t, err)

	assert.Equal(t, ".bash_profile", filepath)
}

// TestCurrentContext will test the currentContext function.
func TestCurrentContext(t *testing.T) {
	_, err := currentContext()
	assert.NoError(t, err)
}

// TestEndpointFromContext will test the endpointFromContext function.
func TestEndpointFromContext(t *testing.T) {
	manifestDir, err := filePathInHomeDir(".docker", "contexts", "meta")
	if err != nil {
		t.Fatal(err)
	}
	// To prevent flakey tests, we will test if the file exists and dynamically
	// assert the result depending on the outcomes.
	if _, err = os.Stat(manifestDir); os.IsNotExist(err) {
		_, err = endpointFromContext("")
		assert.NotNil(t, err)
	} else {
		_, err = endpointFromContext("")
		assert.Nil(t, err)
	}
}
