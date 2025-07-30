package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

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
		Short: "A cli tool to set the default browser on macOS",
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

	if os.Geteuid() == 0 {
		if targetUser == "" {
			return fmt.Errorf("--user must be specified when running as root")
		}
	
		// Look up the specified user and construct the correct plist path
		u, err := user.Lookup(targetUser)
		if err != nil {
			return fmt.Errorf("unknown user %s", targetUser)
		}

		plistPath := filepath.Join(u.HomeDir, "Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist")
		opts = append(opts, client.WithCurrentUser(targetUser), client.WithPlistLocation(plistPath))
	} else {
		if targetUser != "" {
			return fmt.Errorf("--user can only be used when running as root")
		}
	}

	c, err := client.NewClient(opts...)
	if err != nil {
		return err
	}

	err = launchservices.ModifyLS(c, identifier, noRescanLaunchServices)
	if err != nil {
		return err
	}

	// Fix ownership if specifying --user
	if os.Geteuid() == 0 {
		uid, err := parseUID(u)
		if err != nil {
			return err
		}

		err = os.Chown(c.PlistLocation, uid, 20) // gid 20 = staff
		if err != nil {
			return fmt.Errorf("failed to chown plist: %v", err)
		}
	}

	return nil
}

func parseUID(u *user.User) (int, error) {
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return 0, fmt.Errorf("invalid UID for user %s: %v", u.Username, err)
	}
	return uid, nil
}
