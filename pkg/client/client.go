package client

import (
	"os/user"

	osq "github.com/macadmins/osquery-extension/pkg/utils"
)

const launchServicesPlistPath = "Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist"

type Client struct {
	Runner        osq.CmdRunner
	CurrentUser   string
	HomeDir       string
	PlistLocation string
}

type Option func(*Client)
type currentUserLookup func() (*user.User, error)

func WithCurrentUser(currentUser string) Option {
	return func(c *Client) {
		c.CurrentUser = currentUser
	}
}

func WithHomeDir(homeDir string) Option {
	return func(c *Client) {
		c.HomeDir = homeDir
	}
}

func WithPlistLocation(plistLocation string) Option {
	return func(c *Client) {
		c.PlistLocation = plistLocation
	}
}

func NewClient(opts ...Option) (Client, error) {
	return newClient(user.Current, opts...)
}

func newClient(lookupCurrentUser currentUserLookup, opts ...Option) (Client, error) {
	c := Client{}
	c.Runner = osq.NewRunner().Runner
	for _, opt := range opts {
		opt(&c)
	}

	if c.CurrentUser == "" || c.HomeDir == "" {
		currentUser, err := lookupCurrentUser()
		if err != nil {
			return c, err
		}
		if c.CurrentUser == "" {
			c.CurrentUser = currentUser.Username
		}
		if c.HomeDir == "" {
			c.HomeDir = currentUser.HomeDir
		}
	}

	if c.PlistLocation == "" {
		c.PlistLocation = defaultLaunchServicesPlistLocation(c.HomeDir)
	}

	return c, nil
}

func defaultLaunchServicesPlistLocation(homeDir string) string {
	return homeDir + "/" + launchServicesPlistPath
}
