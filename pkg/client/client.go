package client

import (
	"os/user"

	osq "github.com/macadmins/osquery-extension/pkg/utils"
)

type Client struct {
	Runner        osq.CmdRunner
	CurrentUser   string
	PlistLocation string
}

type Option func(*Client)

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
	c := Client{}
	c.Runner = osq.NewRunner().Runner
	for _, opt := range opts {
		opt(&c)
	}

	if c.CurrentUser == "" {
		currentUser, err := user.Current()
		if err != nil {
			return c, err
		}
		c.CurrentUser = currentUser.Username
	}

	if c.PlistLocation == "" {
		userInfo, err := LookupUserInfo(c.CurrentUser)
		if err != nil {
			return c, err
		}
		c.PlistLocation = userInfo.LaunchServicesPlistPath()
	}

	return c, nil
}
