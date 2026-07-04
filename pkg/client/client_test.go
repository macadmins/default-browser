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
	assert.NotEmpty(t, client.CurrentUser, "CurrentUser should not be empty")
	assert.NotEmpty(t, client.HomeDir, "HomeDir should not be empty")
	assert.Equal(t, defaultLaunchServicesPlistLocation(client.HomeDir), client.PlistLocation, "PlistLocation should use the current user's actual home directory")
}

func TestNewClientWithCurrentUser(t *testing.T) {
	expectedUser := "testuser"
	expectedHomeDir := "/Volumes/test_volume/Users/systemuser"

	client, err := newClient(
		func() (*user.User, error) {
			return &user.User{Username: "systemuser", HomeDir: expectedHomeDir}, nil
		},
		WithCurrentUser(expectedUser),
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedUser, client.CurrentUser, "CurrentUser should be set to the provided value")
	assert.Equal(t, expectedHomeDir, client.HomeDir, "HomeDir should be set from the current user lookup")
	assert.Equal(t, defaultLaunchServicesPlistLocation(expectedHomeDir), client.PlistLocation, "PlistLocation should use the current user's actual home directory")
}

func TestNewClientWithHomeDir(t *testing.T) {
	expectedHomeDir := "/Volumes/test_volume/Users/testuser"

	client, err := newClient(
		func() (*user.User, error) {
			return &user.User{Username: "systemuser", HomeDir: "/Users/systemuser"}, nil
		},
		WithHomeDir(expectedHomeDir),
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, "systemuser", client.CurrentUser, "CurrentUser should be set from the current user lookup")
	assert.Equal(t, expectedHomeDir, client.HomeDir, "HomeDir should be set to the provided value")
	assert.Equal(t, defaultLaunchServicesPlistLocation(expectedHomeDir), client.PlistLocation, "PlistLocation should use the provided home directory")
}

func TestNewClientDefaultCurrentUser(t *testing.T) {
	expectedUser := "systemuser"
	expectedHomeDir := "/Volumes/test_volume/Users/systemuser"

	client, err := newClient(
		func() (*user.User, error) {
			return &user.User{Username: expectedUser, HomeDir: expectedHomeDir}, nil
		},
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedUser, client.CurrentUser, "CurrentUser should be set to the system's current user")
	assert.Equal(t, expectedHomeDir, client.HomeDir, "HomeDir should be set to the system's current user's home directory")
	assert.Equal(t, defaultLaunchServicesPlistLocation(expectedHomeDir), client.PlistLocation, "PlistLocation should use the current user's actual home directory")
}

func TestNewClientWithPlistLocation(t *testing.T) {
	expectedPlistLocation := "/tmp/test.plist"

	client, err := newClient(
		func() (*user.User, error) {
			return &user.User{Username: "systemuser", HomeDir: "/Volumes/test_volume/Users/systemuser"}, nil
		},
		WithPlistLocation(expectedPlistLocation),
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedPlistLocation, client.PlistLocation, "PlistLocation should be set to the provided value")
}

func TestNewClientWithCurrentUserAndHomeDirSkipsCurrentUserLookup(t *testing.T) {
	expectedUser := "testuser"
	expectedHomeDir := "/Volumes/test_volume/Users/testuser"

	client, err := newClient(
		func() (*user.User, error) {
			return nil, errors.New("current user lookup should not be called")
		},
		WithCurrentUser(expectedUser),
		WithHomeDir(expectedHomeDir),
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedUser, client.CurrentUser, "CurrentUser should be set to the provided value")
	assert.Equal(t, expectedHomeDir, client.HomeDir, "HomeDir should be set to the provided value")
	assert.Equal(t, defaultLaunchServicesPlistLocation(expectedHomeDir), client.PlistLocation, "PlistLocation should use the provided home directory")
}

func TestNewClientWithCurrentUserHomeDirAndPlistLocationSkipsCurrentUserLookup(t *testing.T) {
	expectedUser := "testuser"
	expectedHomeDir := "/Volumes/test_volume/Users/testuser"
	expectedPlistLocation := "/tmp/test.plist"

	client, err := newClient(
		func() (*user.User, error) {
			return nil, errors.New("current user lookup should not be called")
		},
		WithCurrentUser(expectedUser),
		WithHomeDir(expectedHomeDir),
		WithPlistLocation(expectedPlistLocation),
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedUser, client.CurrentUser, "CurrentUser should be set to the provided value")
	assert.Equal(t, expectedHomeDir, client.HomeDir, "HomeDir should be set to the provided value")
	assert.Equal(t, expectedPlistLocation, client.PlistLocation, "PlistLocation should be set to the provided value")
}

func TestNewClientReturnsCurrentUserLookupError(t *testing.T) {
	expectedErr := errors.New("current user lookup failed")

	client, err := newClient(
		func() (*user.User, error) {
			return nil, expectedErr
		},
	)

	assert.ErrorIs(t, err, expectedErr, "NewClient should return the current user lookup error")
	assert.Empty(t, client.CurrentUser, "CurrentUser should not be set when lookup fails")
	assert.Empty(t, client.HomeDir, "HomeDir should not be set when lookup fails")
	assert.Empty(t, client.PlistLocation, "PlistLocation should not be set when lookup fails")
}

func TestDefaultLaunchServicesPlistLocation(t *testing.T) {
	homeDir := "/Volumes/test_volume/Users/testuser"

	location := defaultLaunchServicesPlistLocation(homeDir)

	assert.Equal(t, homeDir+"/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist", location, "default LaunchServices plist location should use the provided home directory")
}
