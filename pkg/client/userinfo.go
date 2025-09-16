package client

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

type UserInfo struct {
	Username string
	UID      int
	HomeDir  string
}

func (u *UserInfo) LaunchServicesPlistPath() string {
	return filepath.Join(u.HomeDir, "Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist")
}

func LookupUserInfo(username string) (*UserInfo, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, fmt.Errorf("unknown user %s", username)
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return nil, fmt.Errorf("invalid UID for user %s: %v", u.Username, err)
	}

	return &UserInfo{
		Username: u.Username,
		UID:      uid,
		HomeDir:  u.HomeDir,
	}, nil
}

func FixPlistOwnership(username, plistPath string) error {
	userInfo, err := LookupUserInfo(username)
	if err != nil {
		return err
	}

	// Use default group staff (GID 20)
	const staffGID = 20

	if err := os.Chown(plistPath, userInfo.UID, staffGID); err != nil {
		return fmt.Errorf("failed to chown plist: %v", err)
	}

	if err := os.Chmod(plistPath, 0644); err != nil {
		return fmt.Errorf("failed to chmod plist: %v", err)
	}

	return nil
}