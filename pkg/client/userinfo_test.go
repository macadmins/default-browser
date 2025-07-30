package client_test

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/macadmins/default-browser/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestLookupUserInfo(t *testing.T) {
	currentUser, err := user.Current()
	assert.NoError(t, err, "user.Current should not return an error")

	info, err := client.LookupUserInfo(currentUser.Username)
	assert.NoError(t, err, "LookupUserInfo should not return an error")
	assert.Equal(t, currentUser.Username, info.Username, "Username should match")
	assert.Equal(t, currentUser.HomeDir, info.HomeDir, "HomeDir should match")
	assert.Greater(t, info.UID, 0, "UID should be greater than 0")
}

func TestLaunchServicesPlistPath(t *testing.T) {
	userInfo := &client.UserInfo{
		Username: "fakeuser",
		UID:      501,
		HomeDir:  "/Users/fakeuser",
	}

	expected := "/Users/fakeuser/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist"
	assert.Equal(t, expected, userInfo.LaunchServicesPlistPath(), "LaunchServicesPlistPath should construct the correct path")
}

func TestLookupUserInfo_InvalidUser(t *testing.T) {
	_, err := client.LookupUserInfo("fakeuser")
	assert.Error(t, err, "LookupUserInfo should return an error for a non-existent user")
}

func TestFixPlistOwnership(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("TestFixPlistOwnership requires root privileges")
	}

	currentUser, err := user.Current()
	assert.NoError(t, err, "user.Current should not return an error")

	tmpfile, err := os.CreateTemp("", "test.plist")
	assert.NoError(t, err, "CreateTemp should not return an error")
	defer os.Remove(tmpfile.Name())

	err = client.FixPlistOwnership(currentUser.Username, tmpfile.Name())
	assert.NoError(t, err, "FixPlistOwnership should not return an error")

	info, err := os.Stat(tmpfile.Name())
	assert.NoError(t, err, "Stat should not return an error")

	stat := info.Sys().(*os.FileStat)
	assert.Equal(t, uint32(0644), info.Mode().Perm(), "File mode should be 0644")

	// Convert UID to string for comparison with currentUser.Uid
	fileUID := strconv.Itoa(int(stat.Uid))
	assert.Equal(t, currentUser.Uid, fileUID, "File owner UID should match current user")
}
