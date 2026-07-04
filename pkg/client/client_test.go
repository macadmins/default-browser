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
	assert.NotEmpty(t, client.HomeDir, "HomeDir should not be empty")
	assert.Equal(t, defaultLaunchServicesPlistLocation(client.HomeDir), client.PlistLocation, "PlistLocation should use the current user's actual home directory")
}

func TestNewClientWithHomeDir(t *testing.T) {
	expectedHomeDir := "/Volumes/test_volume/Users/testuser"
	expectedPlistLocation := "/Volumes/test_volume/Users/testuser/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist"

	client, err := newClient(
		func() (*user.User, error) {
			return nil, errors.New("current user lookup should not be called")
		},
		WithHomeDir(expectedHomeDir),
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedHomeDir, client.HomeDir, "HomeDir should be set to the provided value")
	assert.Equal(t, expectedPlistLocation, client.PlistLocation, "PlistLocation should use the provided home directory")
}

func TestNewClientDefaultHomeDir(t *testing.T) {
	expectedHomeDir := "/Volumes/test_volume/Users/systemuser"

	client, err := newClient(
		func() (*user.User, error) {
			return &user.User{HomeDir: expectedHomeDir}, nil
		},
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Equal(t, expectedHomeDir, client.HomeDir, "HomeDir should be set to the system's current user's home directory")
	assert.Equal(t, defaultLaunchServicesPlistLocation(expectedHomeDir), client.PlistLocation, "PlistLocation should use the current user's actual home directory")
}

func TestNewClientWithPlistLocation(t *testing.T) {
	expectedPlistLocation := "/tmp/test.plist"

	client, err := newClient(
		func() (*user.User, error) {
			return nil, errors.New("current user lookup should not be called")
		},
		WithPlistLocation(expectedPlistLocation),
	)

	assert.NoError(t, err, "NewClient should not return an error")
	assert.Empty(t, client.HomeDir, "HomeDir should not be set when only PlistLocation is provided")
	assert.Equal(t, expectedPlistLocation, client.PlistLocation, "PlistLocation should be set to the provided value")
}

func TestNewClientWithHomeDirAndPlistLocation(t *testing.T) {
	expectedHomeDir := "/Volumes/test_volume/Users/testuser"
	expectedPlistLocation := "/tmp/test.plist"

	client, err := newClient(
		func() (*user.User, error) {
			return nil, errors.New("current user lookup should not be called")
		},
		WithHomeDir(expectedHomeDir),
		WithPlistLocation(expectedPlistLocation),
	)

	assert.NoError(t, err, "NewClient should not return an error")
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
	assert.Empty(t, client.HomeDir, "HomeDir should not be set when lookup fails")
	assert.Empty(t, client.PlistLocation, "PlistLocation should not be set when lookup fails")
}

func TestDefaultLaunchServicesPlistLocation(t *testing.T) {
	homeDir := "/Volumes/test_volume/Users/testuser"

	location := defaultLaunchServicesPlistLocation(homeDir)

	assert.Equal(t, homeDir+"/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist", location, "default LaunchServices plist location should use the provided home directory")
}
