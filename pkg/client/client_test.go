package client

import (
	"errors"
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
	assert.Equal(t, "/Users/testuser/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist", client.PlistLocation, "PlistLocation should use the provided current user")
}

func TestNewClientDefaultCurrentUser(t *testing.T) {
	client, err := newClient(func() (*user.User, error) {
		return &user.User{Username: "systemuser"}, nil
	})
	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, "systemuser", client.CurrentUser, "CurrentUser should be set to the system's current user")
	assert.Equal(t, "/Users/systemuser/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist", client.PlistLocation, "PlistLocation should use the system's current user")
}

func TestNewClientWithPlistLocation(t *testing.T) {
	expectedPlistLocation := "/tmp/test.plist"
	client, err := NewClient(WithPlistLocation(expectedPlistLocation))
	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedPlistLocation, client.PlistLocation, "PlistLocation should be set to the provided value")
}

func TestNewClientWithCurrentUserSkipsCurrentUserLookup(t *testing.T) {
	client, err := newClient(func() (*user.User, error) {
		return nil, errors.New("current user lookup should not be called")
	}, WithCurrentUser("testuser"))

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, "testuser", client.CurrentUser, "CurrentUser should be set to the provided value")
	assert.Equal(t, "/Users/testuser/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist", client.PlistLocation, "PlistLocation should use the provided current user")
}

func TestNewClientReturnsCurrentUserLookupError(t *testing.T) {
	expectedErr := errors.New("current user lookup failed")

	client, err := newClient(func() (*user.User, error) {
		return nil, expectedErr
	})

	assert.ErrorIs(t, err, expectedErr, "NewClient should return the current user lookup error")
	assert.Empty(t, client.CurrentUser, "CurrentUser should not be set when lookup fails")
}

func TestDefaultLaunchServicesPlistLocation(t *testing.T) {
	location := defaultLaunchServicesPlistLocation("testuser")

	assert.Equal(t, "/Users/testuser/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist", location, "default LaunchServices plist location should match the legacy path")
}
