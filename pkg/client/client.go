package client

import (
	"os/user"

	osq "github.com/macadmins/osquery-extension/pkg/utils"
)

const launchServicesPlistPath = "Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist"

type Client struct {
	Runner        osq.CmdRunner
	HomeDir       string
	PlistLocation string
}

type Option func(*Client)
type currentUserLookup func() (*user.User, error)

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

	if c.HomeDir == "" && c.PlistLocation == "" {
		currentUser, err := lookupCurrentUser()
		if err != nil {
			return c, err
		}
		c.HomeDir = currentUser.HomeDir
	}

	if c.PlistLocation == "" {
		c.PlistLocation = defaultLaunchServicesPlistLocation(c.HomeDir)
	}

	return c, nil
}

func defaultLaunchServicesPlistLocation(homeDir string) string {
	return homeDir + "/" + launchServicesPlistPath
}
