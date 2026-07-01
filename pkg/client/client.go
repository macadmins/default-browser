package client

import (
	"os"
	"os/user"

	osq "github.com/macadmins/osquery-extension/pkg/utils"
)

const launchServicesPlistPath = "Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist"

type Client struct {
	Runner        osq.CmdRunner
	CurrentUser   string
	PlistLocation string
}

type Option func(*Client)
type currentUserLookup func() (*user.User, error)
type homeDirLookup func() (string, error)

func WithCurrentUser(currentUser string) Option {
	return func(c *Client) {
		c.CurrentUser = currentUser
	}
}

func WithPlistLocation(plistLocation string) Option {
	return func(c *Client) {
		c.PlistLocation = plistLocation
	}
}

func NewClient(opts ...Option) (Client, error) {
	return newClient(user.Current, os.UserHomeDir, opts...)
}

func newClient(lookupCurrentUser currentUserLookup, lookupHomeDir homeDirLookup, opts ...Option) (Client, error) {
	c := Client{}
	c.Runner = osq.NewRunner().Runner
	for _, opt := range opts {
		opt(&c)
	}

	if c.CurrentUser == "" {
		currentUser, err := lookupCurrentUser()
		if err != nil {
			return c, err
		}
		c.CurrentUser = currentUser.Username
	}

	if c.PlistLocation == "" {
		homeDir, err := lookupHomeDir()
		if err != nil {
			return c, err
		}
		c.PlistLocation = defaultLaunchServicesPlistLocation(homeDir)
	}

	return c, nil
}

func defaultLaunchServicesPlistLocation(homeDir string) string {
	return homeDir + "/" + launchServicesPlistPath
}
