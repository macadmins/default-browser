package client

import (
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
	return newClient(user.Current, opts...)
}

func newClient(lookupCurrentUser currentUserLookup, opts ...Option) (Client, error) {
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
		c.PlistLocation = defaultLaunchServicesPlistLocation(c.CurrentUser)
	}

	return c, nil
}

func defaultLaunchServicesPlistLocation(currentUser string) string {
	return "/Users/" + currentUser + "/" + launchServicesPlistPath
}
