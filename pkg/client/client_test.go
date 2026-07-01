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
	expectedHomeDir := "/Volumes/test_volume/Users/testuser"

	client, err := newClient(
		func() (*user.User, error) {
			return nil, errors.New("current user lookup should not be called")
		},
		func() (string, error) {
			return expectedHomeDir, nil
		},
		WithCurrentUser(expectedUser),
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedUser, client.CurrentUser, "CurrentUser should be set to the provided value")
	assert.Equal(t, expectedHomeDir+"/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist", client.PlistLocation, "PlistLocation should use the current user's actual home directory")
}

func TestNewClientDefaultCurrentUser(t *testing.T) {
	expectedHomeDir := "/Volumes/test_volume/Users/systemuser"

	client, err := newClient(
		func() (*user.User, error) {
			return &user.User{Username: "systemuser"}, nil
		},
		func() (string, error) {
			return expectedHomeDir, nil
		},
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, "systemuser", client.CurrentUser, "CurrentUser should be set to the system's current user")
	assert.Equal(t, expectedHomeDir+"/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist", client.PlistLocation, "PlistLocation should use the current user's actual home directory")
}

func TestNewClientWithPlistLocation(t *testing.T) {
	expectedPlistLocation := "/tmp/test.plist"

	client, err := NewClient(WithPlistLocation(expectedPlistLocation))

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedPlistLocation, client.PlistLocation, "PlistLocation should be set to the provided value")
}

func TestNewClientWithCurrentUserSkipsCurrentUserLookup(t *testing.T) {
	expectedHomeDir := "/Volumes/test_volume/Users/testuser"

	client, err := newClient(
		func() (*user.User, error) {
			return nil, errors.New("current user lookup should not be called")
		},
		func() (string, error) {
			return expectedHomeDir, nil
		},
		WithCurrentUser("testuser"),
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, "testuser", client.CurrentUser, "CurrentUser should be set to the provided value")
	assert.Equal(t, expectedHomeDir+"/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist", client.PlistLocation, "PlistLocation should use the current user's actual home directory")
}

func TestNewClientReturnsCurrentUserLookupError(t *testing.T) {
	expectedErr := errors.New("current user lookup failed")

	client, err := newClient(
		func() (*user.User, error) {
			return nil, expectedErr
		},
		func() (string, error) {
			return "/Volumes/test_volume/Users/testuser", nil
		},
	)

	assert.ErrorIs(t, err, expectedErr, "NewClient should return the current user lookup error")
	assert.Empty(t, client.CurrentUser, "CurrentUser should not be set when lookup fails")
}

func TestNewClientReturnsHomeDirLookupError(t *testing.T) {
	expectedErr := errors.New("home directory lookup failed")

	client, err := newClient(
		func() (*user.User, error) {
			return &user.User{Username: "testuser"}, nil
		},
		func() (string, error) {
			return "", expectedErr
		},
	)

	assert.ErrorIs(t, err, expectedErr, "NewClient should return the home directory lookup error")
	assert.Equal(t, "testuser", client.CurrentUser, "CurrentUser should be set before home directory lookup")
	assert.Empty(t, client.PlistLocation, "PlistLocation should not be set when home directory lookup fails")
}

func TestDefaultLaunchServicesPlistLocation(t *testing.T) {
	homeDir := "/Volumes/test_volume/Users/testuser"

	location := defaultLaunchServicesPlistLocation(homeDir)

	assert.Equal(t, homeDir+"/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist", location, "default LaunchServices plist location should use the provided home directory")
}
