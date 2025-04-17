package launchservices

import (
	"os"
	"testing"

	"github.com/macadmins/default-browser/pkg/client"
	"github.com/macadmins/osquery-extension/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestRebuildLaunchServices(t *testing.T) {
	mockRunner := utils.MockCmdRunner{
		Output: "",
		Err:    nil,
	}

	mockClient := client.Client{
		Runner: mockRunner,
	}

	err := rebuildLaunchServices(mockClient, false)
	assert.NoError(t, err, "rebuildLaunchServices should not return an error")

	err = rebuildLaunchServices(mockClient, true)
	assert.NoError(t, err, "rebuildLaunchServices should not return an error when noRebuildLaunchServices is true")
}

func TestKilllsd(t *testing.T) {
	mockRunner := utils.MockCmdRunner{
		Output: "",
		Err:    nil,
	}

	mockClient := client.Client{
		Runner: mockRunner,
	}

	err := killlsd(mockClient)
	assert.NoError(t, err, "killlsd should not return an error")
}

func TestLsRegisterExists(t *testing.T) {
	// Create a temporary file to test existence
	tmpfile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err, "CreateTemp should not return an error")
	defer os.Remove(tmpfile.Name())

	// Test that the file exists
	exists, err := lsRegisterExists(tmpfile.Name())
	assert.NoError(t, err, "lsRegisterExists should not return an error")
	assert.True(t, exists, "lsRegisterExists should return true for an existing file")

	// Remove the temporary file
	err = os.Remove(tmpfile.Name())
	assert.NoError(t, err, "Remove should not return an error")

	// Test that the file does not exist
	exists, err = lsRegisterExists(tmpfile.Name())
	assert.NoError(t, err, "lsRegisterExists should not return an error for a non-existing file")
	assert.False(t, exists, "lsRegisterExists should return false for a non-existing file")
}
