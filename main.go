package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/macadmins/default-browser/pkg/client"
	"github.com/macadmins/default-browser/pkg/launchservices"
	"github.com/spf13/cobra"
)

var version string

func main() {
	var identifier string
	var noRescanLaunchServices bool
	var targetUser string

	var rootCmd = &cobra.Command{
		Use:   "default-browser",
		Short: "A CLI tool to set the default browser on macOS",
		RunE: func(cmd *cobra.Command, args []string) error {
			return setDefault(identifier, noRescanLaunchServices, targetUser)
		},
	}

	rootCmd.Flags().StringVar(&identifier, "identifier", "", "An identifier for the application")
	rootCmd.Flags().BoolVar(&noRescanLaunchServices, "no-rescan-launchservices", false, "Do not rescan launch services.")
	rootCmd.Flags().BoolVar(&noRescanLaunchServices, "no-rebuild-launchservices", false, "Legacy: same as --no-rescan-launchservices")
	rootCmd.Flags().StringVar(&targetUser, "user", "", "Username to operate on (only allowed when run as root)")

	rootCmd.MarkFlagRequired("identifier")
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("default-browser version {{.Version}}\n")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func setDefault(identifier string, noRescanLaunchServices bool, targetUser string) error {
	var opts []client.Option
	var plistPath string

	if os.Geteuid() == 0 {
		if targetUser == "" {
			return fmt.Errorf("--user must be specified when running as root")
		}

		userInfo, err := client.LookupUserInfo(targetUser)
		if err != nil {
			return err
		}

		plistPath = filepath.Join(userInfo.HomeDir, "Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist")
		opts = append(opts, client.WithCurrentUser(userInfo.Username), client.WithPlistLocation(plistPath))
	} else {
		if targetUser != "" {
			return fmt.Errorf("--user can only be used when running as root")
		}
	}

	c, err := client.NewClient(opts...)
	if err != nil {
		return err
	}

	if err := launchservices.ModifyLS(c, identifier, noRescanLaunchServices); err != nil {
		return err
	}

	if os.Geteuid() == 0 {
		if err := client.FixPlistOwnership(targetUser, c.PlistLocation); err != nil {
			return err
		}
	}

	return nil
}
