package client

import (
	"os/user"
	"testing"

	osq "github.com/macadmins/osquery-extension/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	assert.NoError(t, err, "NewClient should not return an error")
	assert.NotNil(t, client.Runner, "Runner should not be nil")
	assert.IsType(t, &osq.ExecCmdRunner{}, client.Runner, "Runner should be of type osq.Runner")
}

func TestNewClientWithCurrentUser(t *testing.T) {
	expectedUser := "testuser"
	client, err := NewClient(WithCurrentUser(expectedUser))
	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedUser, client.CurrentUser, "CurrentUser should be set to the provided value")
}

func TestNewClientDefaultCurrentUser(t *testing.T) {
	currentUser, err := user.Current()
	assert.NoError(t, err, "user.Current should not return an error")

	client, err := NewClient()
	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, currentUser.Username, client.CurrentUser, "CurrentUser should be set to the system's current user")
}

func TestNewClientWithPlistLocation(t *testing.T) {
	expectedPlistLocation := "/tmp/test.plist"
	client, err := NewClient(WithPlistLocation(expectedPlistLocation))
	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedPlistLocation, client.PlistLocation, "PlistLocation should be set to the provided value")
}
